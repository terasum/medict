# Changelog

All notable changes to [Medict](https://github.com/terasum/medict) are recorded
here. The most recent release is at the top. For the full history before
v3.1.0, see `git log v3.0.1..HEAD`.

## v3.1.4

Bookmarks / saved-words feature + navigation restructure.

### Added
- **Bookmarks / 生词本** (#643): save words while looking them up (star button in
  the search bar), then browse and re-lookup them in a dedicated **「生词」tab**.
  - Star button next to back/forward in the search bar — saves the current word
    + dictionary context.
  - Full-page bookmarks list with search/filter, click-to-lookup, and delete.
  - Persisted to `bookmarks.json` in the app config directory.

### Changed
- **Navigation restructure**: the left sidebar tabs are now **搜索 / 生词 / 词典 /
  设置**. Plugins and Debug moved into the Settings page as sub-navigation items
  (插件设置 already existed; 调试工具 newly added).

## v3.1.3

New features + bug fixes from the stale-issue review (#258/#259/#260), plus the
TypeScript strict-mode type-safety gate (#702) and backend/frontend cleanup.

### Added
- **Double-click to look up a word** (#258): double-clicking text in a
  dictionary entry selects it (using the browser's CJK-aware word-boundary
  selection) and looks it up — matching macOS Dictionary.app behavior.
- **Keyboard & wheel zoom** (#260): `Ctrl/Cmd + =/-` and `Ctrl/Cmd + scroll`
  now zoom the definition iframe (previously only toolbar buttons).

### Fixed
- **`@@@LINK` multi-target & chained redirects** (#260): a headword redirecting
  to multiple targets (e.g. 滋 → 滋 + 滋) no longer shows blank; chained
  redirects (A→B→C) now follow through instead of stopping at the first hop or
  leaking raw `@@@LINK=` text. Cycle-guarded, max depth 8.
- **Dictionary name truncation** (#259): `mdict.Name()` used `TrimRight` (cutset)
  which ate trailing `m`/`d`/`x` chars — `mad.mdx` became `ma`. Fixed to
  `TrimSuffix` of the actual extension.

### Changed
- **TypeScript strict mode** (#702): `tsconfig.json` now has `strict: true`;
  `vue-tsc --noEmit` passes with zero errors and runs as a CI gate. All 66
  strict errors fixed (typed `queryPendingList`/`selectDict`, debounce params,
  `$onAction` callback, `*.vue` shim, etc.).
- **Backend naming consistency** (#737): `back_server.go` + `back_server_inner.go`
  merged into one file; `DictCon → Controller`, `MainDict → Dict`, `dc.ds → svc`.
- **Startup logging** (#735): `[medict-init]` bootstrap `fmt.Printf` → logger.

## v3.1.2

Maintenance release: internal architecture cleanup — no user-visible behavior
change. Three refactors land together (all backend builds + frontend `bun build`
green; CI passing).

### Changed
- **Typed IPC (#729)**: the frontend↔backend call path no longer goes through a
  string-dispatched `App.Dispatch(apiName, map)` + runtime type assertions. Each
  handler is now a typed `App` method (`InitDicts` / `GetAllDicts` /
  `SearchWord(dictId, word)` / `BuildIndexByDictId(dictid)`), the same proven
  Wails native-binding pattern the app already used for `ResourceServerAddr` etc.
  Removes `Dispatch` / `handlerMap` / `DispatchIPCReq` and three dead
  no-backend frontend calls.
- **`DictService` dependency injection (#727)**: removed the package-global
  singleton (`sync.Once` + `GetDictService()`); `App` now owns the `*DictService`
  and injects it into `BackServer.SetUp`. Improves testability and constructor
  injection.
- **Startup/query log trimming (#735)**: per-Locate/Lookup `INFO` logs in go-mdict
  and the holder downgraded to `DEBUG`; two stray `fmt.Printf` removed. Reduces
  log I/O and noise during startup/indexing. (BuildIndex's per-keyword path was
  already log-free.)

> Note: typed IPC changes the frontend↔backend binding contract internally; the
> typed-binding mechanism is already in use app-wide, so runtime behavior is
> unchanged. `model.Resp.data` typing (`as unknown as` casts) is unaffected —
> that's tracked separately in #702.

## v3.1.1

Patch release: fixes a **live crash** during index build, `entry://` link
navigation, and a class of keyword-index correctness bugs from the indexer
review (#722), plus architecture cleanups. No data-format breakage (existing
`.melev` caches are rebuilt once on the new schema).

### Fixed
- **Live panic on index build (#724)**: the leveldb connection pool could hand
  out a nil connection and panic during `BuildIndex`. Pool removed; a single
  shared `*leveldb.DB` handle is used (goleveldb is concurrency-safe).
- **`entry://` links with `=`/`/`/`.` (#718/#723)**: no longer trigger an OS
  "There is no application set to open the URL" dialog; in-iframe jumps work.
- **Prefix search for P/F/K/W-initial words (#723)**: a `TrimLeft`-vs-`TrimPrefix`
  bug had collapsed distinct words onto one key and broken prefix search.
- **Same-headword entries / homographs (#742)**: the leveldb key now carries the
  record offset, so "bank" sense 1 & 2 no longer overwrite each other — both
  surface in search. (One-time index rebuild via `schema_version` bump.)
- **Dictionary list ordering (#740)**: sorted by name, not by dir-path MD5.
- **stardict (#739)**: `Lookup`/`LookupResource` return a clean `ErrNotFound`
  instead of swallowing errors / returning empty 200s.
- macOS release now ships a `.app.zip` alongside the `.dmg` (#721).

### Changed
- Indexer `Close` lifecycle end-to-end (LvDB → indexer → holder → service →
  `App.shutdown`); the `.mdx` parser opens/closes per op, so no handle leak (#731).
- Index built via a single leveldb **batch** write; the fuzzy BK-tree is built
  from the leveldb index (no `.mdx` re-parse on cache hit) (#741).
- `model.ErrNotFound` sentinel — a normal miss is distinguishable from a real
  error via `errors.Is` (#732).
- Word lookup is now an **explicit Gin route** (`GET /__mdict/__tcidem_query`);
  resource lookups stay a catch-all (#728).
- Boilerplate reduction, dead-code removal, log consolidation (#734/#735/#736/#737/#738).

## v3.1.0

First release since v3.0.1. Bundles the reliability, fuzzy-search, and build
modernization work from the 2024–2026 cycle.

### Added

- **Mdict fuzzy search** (#666 / #697): lookups are prefix-first, with an
  in-memory BK-tree (character-level Levenshtein, tolerance 2) fallback when the
  prefix misses — so near-misses still find the entry.
- **leveldb indexer** is now the active index engine; each `.mdx` gets a
  `.melev` sidecar cache, and the BK-tree is rebuilt lazily for users whose
  cache skips `BuildIndex`.
- **PR-level CI** (backend `go build`/`vet`/`test` + frontend `bun build`)
  on every pull request to `develop`/`master`. (#701)
- **Architecture guide** (`CLAUDE.md`) documenting the two-channel (Wails IPC +
  embedded Gin server) design for contributors and AI assistants.

### Fixed

- `entry://` links whose IDs contain `=`, `/`, `.` etc. (e.g.
  `entry://topic_transport-by-water_level=c1`) no longer trigger a system
  "There is no application set to open the URL" dialog — in-iframe jumps now
  work, with a client-side click interceptor as a safety net. (#718)
- **Word-lookup main flow** (#707): declarative iframe, fixed a global
  keyboard-listener leak, and a last-write-wins guard on `searchWord` so fast
  consecutive lookups no longer race.
- **Startup robustness** (#705): init failures surface as a dialog instead of
  crashing the app.
- **History stack back/forward** (#699): a `keyword`/`key_word` field mix-up
  broke navigation; the field is now unified to `keyword`.
- **Wails IPC browser fallback** (#700): the browser path returned `undefined`
  due to a missing `return`.
- `GracefulStop` no longer calls `log.Fatal`/`os.Exit` on a shutdown error. (#683)
- `DictService` singleton is now race-free (`sync.Once`).
- **CORS** restricted to a local allowlist instead of reflecting any `Origin`.
- Correct `Content-Type` for `woff2` fonts; correct static-server URL prefix
  check.
- Stardict index type distinguished from Mdict; symlinked dictionary
  directories are followed.

### Changed

- Frontend package manager migrated **pnpm → bun**. (#695)
- Backend loggers given meaningful per-package identifiers (#704);
  `DictService` lock granularity tightened so no I/O runs while holding the
  lock.
- Frontend dead code and unused dependencies removed. (#703)

### Build / Infrastructure

- Bumped to **Go 1.25**; upgraded `golang.org/x/*` and `logrus`. (#690)
- Re-enabled and modernized the **release pipeline**: tag-triggered, Go 1.25,
  modern actions, three platforms (macOS universal, Linux x86_64, Windows
  x86_64), correct version injection on every platform.
- Cleared `develop`-branch `go vet`/test debt to keep CI green. (#715)

> Note: the previous automated release workflow had been silently failing on
> every `develop` push for months (stale Go version, retired runners, deprecated
> actions). It has been replaced by the tag-triggered pipeline above.
