# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Medict v3 is a cross-platform desktop dictionary app built with **Wails v2** (Go 1.25 backend + Vue 3/Vite frontend, GPL-3.0). It queries `.mdx/.mdd` (Mdict v1.x & v2.0) and stardict dictionaries. Module path: `github.com/terasum/medict`.

## Common commands

Development (requires the Wails CLI installed: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`):
- `make dev` → `wails dev` — runs backend + frontend together with hot reload (the canonical dev loop)
- `make build` → `wails build -devtools` — produces a binary in `build/bin/`
- `make license` → re-applies GPL headers via `addlicense`

Go tests (standard module):
- All: `go test ./...`
- One package: `go test ./pkg/service/mdict/...`
- One test: `go test ./pkg/service/mdict/ -run TestMdictHolder -v`

Frontend lives in `frontend/`. Package manager is **bun** (migrated from pnpm via #695; `wails.json` and CI use bun):
- `cd frontend && bun install`
- `bun run dev` / `bun run build`

Linux build additionally needs `libgtk-3-dev libwebkit2gtk-4.0-dev`.

## Architecture

### Two communication channels to the frontend
This is the single most important thing to understand. The Go side talks to the Vue webview over **two independent paths**:

1. **Wails IPC (control plane)** — `App` methods in `app.go` are bound to JS (bindings are generated into `frontend/wailsjs/`). `App.Dispatch(apiName, args)` routes to a `handlerMap` of named handlers (`InitDicts`, `GetAllDicts`, `SearchWord`, `BuildIndexByDictId`) registered in `pkg/backserver/back_server.go` `setupHandlers()`.

2. **Embedded Gin HTTP server (data plane)** — dictionary definitions are rendered as full HTML inside an `<iframe>` in the webview. Those HTML payloads, plus their resources (css/js/images/fonts/audio), must be served same-origin, so a **Gin server is started on `localhost:0` (random port)** at app init. The frontend learns the port at runtime via the `App.ResourceServerAddr()` IPC call. URL root is `static.ContentRootUrl = "/__mdict"`. Word lookups use magic path `/__tcidem_query`.

   Word lookups go through an **explicit route** `GET /__mdict/__tcidem_query` (`HandleWordQueryReq`); everything else under the content root is a resource lookup handled by a `NoRoute` catch-all (`HandleResourceQueryReq`) (`back_server_inner.go setUpRouters`). Resource paths are arbitrary (css/images/fonts/...), so they stay a catch-all rather than enumerated routes.

### Dictionary abstraction
`pkg/model/dict_interface.go` defines `GeneralDictionary` (`BuildIndex/Lookup/Locate/Search/LookupResource/DictType`). `mdict` (`pkg/service/mdict`) and `stardict` (`pkg/service/stardict`) implement it. `mdictSvcImpl` wraps one `mdx` holder plus N `mdd` holders (mdds are searched in order for resources).

`pkg/service/dicts_service.go` `DictService` holds `map[string]*model.DictionaryItem`; it's constructed by `App` via `NewDictService` and injected into `BackServer.SetUp` (no longer a package-global singleton — #727). The map key is the MD5 of the dictionary directory path (see `pkg/service/dict_item.go` `NewByDirItem`). Every method takes `dictLock`. Lookup/Search/Locate all **require `BuildIndex` to have run first** or they return `"dictionary not ready"`.

### Dictionaries are directories, auto-scanned
A dictionary = one directory under `BaseDictDir` (config `medict.toml`, default per-OS app-data `.../medict/dicts`). `pkg/service/support/filewalker.go` walks the dir; `DirItem` (`pkg/model/dict_def.go`) captures mdx/mdd or stardict (dz/ifo/idx) files plus optional `_cover.jpg`, `_mdict.dtype`/`_stardict.dtype`. Type is auto-detected by file presence.

### Mdict indexing & search
`internal/libs/go-mdict` parses the binary format. The **active indexer is leveldb** (`pkg/service/mdict/mdict-idxer/leveldb_indexer.go`, via `NewIndexer`); the sqlite indexer in the same dir is dead/test-only. It builds a `.melev` sidecar next to each `.mdx`.

`mdictHolder.BuildIndex` (`pkg/service/mdict/mdict_holder.go`) drives the build: pulls all keywords via `rawdict.GetKeyWordEntries()`, converts each to an offset-bearing `*MdictKeyWordIndex` (`ConvertKeyWordIndex`), and `AddRecord`s it. Idempotency: if meta `entries_num` is set, the build is skipped.

`mdictHolder.Search` is **prefix-first + fuzzy fallback** (#697): the leveldb indexer does a prefix scan; on a miss it falls back to an in-memory **BK-tree** (Levenshtein, tolerance 2, limit 100, sorted by distance). The BK-tree is built eagerly in `BuildIndex` and lazily rebuilt by `ensureBkTree()` on first miss (covers users whose `.melev` cache skips `BuildIndex`). It wraps entries in a `fuzzyEntry` whose `Distance` uses **direct character-level Levenshtein** — deliberately NOT `MdictKeyWordIndex.Distance`, which runs on `utils.StrToUnicode`-escaped strings (every rune → 6 chars, inflating distance 6× and breaking tolerance).

Supporting data-structure libs (ART, patricia trie, bktree) live under `internal/libs/` and are vendored/local — treat them as part of this repo, not external deps.

### Content pipeline (resource URL rewriting)
Mdict HTML references resources by relative paths that only resolve inside the embedded Gin server. `internal/static/handler/` is a pipeline of `replacer_*` passes (css, image, javascript, link, sound, entry) that rewrite those references and inline/transform content. `handler.WrapContent` and `handler.WrapResource` are the entry points called from `pkg/apis/dicts_controller.go`. `entry://` jumps and `@@@LINK=` redirects are handled in `HandleWordQueryReq`.

### Startup sequence
`main.go` → `NewApp()` → `appInit()` synchronously loads config (`internal/entry/app_loader.go` `LoadApp`, which also writes a default `medict.toml` and unpacks a preset dictionary) → builds & starts the Gin server → Wails lifecycle (`startup`/`domReady`/`shutdown`) starts channel listeners. Init errors are pushed to `errorChannel` and surfaced as a dialog.

## Known sharp edges (verify before relying on)

These are non-obvious behaviors that span multiple files:

- **`back_server.go` vs `back_server_inner.go` split isn't obvious from the names**: `StaticServerBaseUrl`/`GracefulStop`/`Start`/`SetUp`/`DispatchIPCReq` live in `back_server.go`; `startStaticServer`/`cors()`/`setUpRouters`/`setupHandlers` live in `back_server_inner.go`.
- **Dev-mode Gin port is `localhost:9081`**, but `SetDebug()` is never called on the current startup path, so production always binds a random port resolved via `ResourceServerAddr()`.
- Frontend `frontend/wailsjs/` is **generated** by Wails from the `App` struct — do not hand-edit; it regenerates on `wails dev`/`wails build`.
