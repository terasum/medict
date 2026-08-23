# Changelog

All notable changes to [Medict](https://github.com/terasum/medict) are recorded
here. The most recent release is at the top. For the full history before
v3.1.0, see `git log v3.0.1..HEAD`.

## v3.1.13

Platform-aware title bar, ECDICT rendering fix, and inflection fallback for the built-in dictionary.

### Added
- **Platform detection** (`App.Platform()`): the frontend can now query the
  host OS (runtime.GOOS) over IPC. On Windows and Linux, where the native
  title bar cannot be hidden, the in-app fake title bar shrinks by 12px
  (26px → 14px), reclaiming vertical space for dictionary content. macOS
  keeps the full-height strip for window dragging.
- **ECDICT inflection fallback** (part of #786): searching or clicking an
  inflected English form now falls back to its base form — `parts` → `part`,
  `studies` → `study`, `boxes` → `box`, `ran` / `went` / `written` via an
  irregular-verb table, `-ing` forms (`lying` → `lie`, `running` → `run`),
  and comparatives (`bigger` → `big`, `happiest` → `happy`). Works in both
  exact lookup and prefix search, including in-entry hyperlink lookups.
  Non-English (e.g. Chinese) input is exempt from the English-only rules.

### Fixed
- **ECDICT line breaks**: ECDICT's source CSV encodes in-field line breaks as
  a literal `\n` (backslash-n), which the importer stores verbatim — so
  multi-sense definitions rendered as one run-on line with visible `\n`
  text (10k+ rows in the bundled preset affected). The renderer now converts
  both the literal `\n` sequence and real newlines to `<br>`. No database
  reinstall is needed.

## v3.1.12

Performance: dictionary BuildIndex 51x faster, ECDICT Search 100x faster, BK-tree persistence, benchmark harness.

### Added
- **Benchmark harness** (#806, #807): `cmd/benchmark` tool profiles both built-in
  dictionaries (ECDICT + cc-cedict) with per-scenario CPU/heap profiling, outputs
  structured JSON for AI analysis.

### Changed
- **ECDICT Search 100x faster** (#808): `PRAGMA case_sensitive_like = ON` lets
  SQLite use the primary key index for `LIKE 'prefix%'` — prefix search on 50K
  entries drops from 12ms to 0.1ms. Regression guard test prevents removal.

### Fixed
- **Dictionary BuildIndex 51x faster** (#781): the BK-tree (fuzzy search
  structure) consumed 92% of BuildIndex time (~2 minutes for 195K entries).
  Now built asynchronously in a background goroutine — BuildIndex drops from
  2m3s to 2.4s. The UI no longer freezes when adding large dictionaries.
- **BK-tree persistence**: the in-memory BK-tree is now serialized to
  `.melev/bktree.gob` after async build. On subsequent launches it loads in
  ~1.6s (vs 2-minute rebuild), making fuzzy search immediately available.
- **Search non-blocking on BK-tree miss**: if the background BK-tree build is
  still running, Search misses skip fuzzy (return empty) instead of blocking.

## v3.1.11

Dictionary groups, persistent definition zoom, a unified settings experience,
and a more reliable offline-first default dictionary.

### Added
- **Dictionary groups**: organize dictionaries from the dictionary sidebar and
  manage groups with the same bottom action-bar pattern used by bookmarks.
- **Definition zoom** (#260): the search-page zoom controls now change the whole
  definition document in 10% steps from 80% to 160%, work in single- and
  multi-dictionary modes, and persist across reloads.
- **Standalone CSS editor** (#783): dictionary CSS editing now opens in its own
  Wails window, starts from the dictionary's existing styles when available,
  and keeps live preview support.
- **Settings documentation tabs**: usage, privacy, and license documents now
  live inside Settings with a stable sidebar and top-level document tabs.

### Changed
- **Offline-first default dictionary**: removed the unreliable Bing default and
  added an in-app installer for the full ECDICT dataset while retaining the
  bundled subset as the always-available fallback.
- **Settings and debug tools**: flattened the settings hierarchy, merged version
  information into About, and embedded the resource-query debugger directly in
  the Proton-styled settings layout.
- **Definition lookup interaction**: hovering marks word boundaries with an
  underline; lookup occurs only after clicking, without a popup.
- **Interface consistency**: unified header alignment, dictionary icons and
  active state, bookmark/dictionary sidebars, action bars, controls, and footer
  placement. Vite was upgraded to version 6.

### Fixed
- Definition-link and sidebar lookups now update both the active entry and the
  search field, and sidebar entries remain clickable after navigation.
- Dictionary-provided cover icons are preferred, with a shared fallback icon
  for dictionaries that do not include one.
- Online dictionary failures degrade cleanly instead of returning unhandled 500
  responses.

## v3.1.10

Hotfix: CSS editor storage moved to app config dir (fixes readonly filesystem error).

### Fixed
- **CSS editor saves to app config dir** (#783): the per-dictionary CSS
  override was stored in the dictionary directory (`_medict_user.css`), which
  fails on read-only filesystems (EROFS). Now stored under
  `<app_config_dir>/user_css/<dictId>.css` (always writable). `WrapContent`
  reads from an in-memory map loaded at startup instead of reading a file per
  render.

## v3.1.9

Dictionary CSS editor (CodeMirror + live preview) + photon framework restore (fixes layout regressions from v3.1.8's UnoCSS migration).

### Added
- **Dictionary CSS editor** (#783): a CodeMirror-powered CSS editor accessible
  from the lookup-page toolbar. Edit per-dictionary CSS overrides with live
  preview (client-side iframe injection, 300ms debounce — no round-trip).
  Overrides persist to a sidecar `_medict_user.css` and auto-inject on every
  future lookup via `WrapContent`.

### Fixed
- **Layout regressions from UnoCSS migration** (v3.1.8): the photon CSS
  framework was retired in v3.1.8 and replaced with UnoCSS shortcuts that
  didn't match photon's exact layout — causing sidebar lists to disappear,
  settings sidebar to go horizontal, main layout overflow, and toolbar button
  misalignment. Photon is now **restored** (the UnoCSS token color unification
  is kept); the conflicting shortcuts removed.

## v3.1.8

Default EN-CN dictionary (offline ECDICT + free online Bing), recursive dict-scan fix, a frontend style unification (UnoCSS), and frontend test infrastructure.

### Added
- **Default English-Chinese dictionary** (#794, #795): Medict now ships with an
  offline EN-CN dictionary out of the box — a curated **ECDICT** subset (~50k
  high-frequency words, embedded, no key, no network) as the default, plus a
  free **Bing online** dictionary (scrapes `cn.bing.com/dict/clientsearch`, no
  API key). Both appear in the dict list and participate in multi-dict /
  hover-popup / Anki export like any other dictionary.
- **Frontend test infrastructure** (#801): introduced **vitest** + happy-dom +
  `@vue/test-utils`; first batch of unit tests covering the 3 Pinia stores
  (ui / bookmark getters / dict `buildEntryURL`).

### Changed
- **Frontend style unification — UnoCSS** (#800): adopted **UnoCSS**
  (presetWind3, Tailwind-compatible) as the atomic CSS layer; established a
  single design-token source (`tokens.ts`: brand `#326cb8`, gray scale,
  danger, font stacks) shared by **both** UnoCSS and the naive-ui
  `themeOverrides` — no more off-brand colors. All components' hardcoded
  colors unified to tokens; the legacy GitHub **photon** CSS framework
  retired (−2352 lines, classes replaced by UnoCSS shortcuts).

### Fixed
- **Recursive dict-scan** (#257, #797): nested dict directories (e.g.
  `English/OALD9/oald9.mdx`) are now found exactly once — category folders are
  no longer mis-attributed as dictionaries, and duplicates are eliminated.

## v3.1.7

Multi-dict query, Anki export, hover-popup lookup, and a bookmarks v2 overhaul.

### Added
- **Multi-dict simultaneous query** (#776, #778): a 「多」toggle in the dict
  toolbar turns the dict icons into multi-select; searching then queries every
  selected dictionary and stacks the definitions vertically (GoldenDict-style).
  The active set is persisted across restarts via the settings write-back (#777).
- **Export bookmarks to Anki** (#774, #775): the 生词本 page can export its words
  to a native `.apkg` — each notebook becomes an Anki deck, and each card carries
  the saved HTML snapshot with images extracted as Anki media (validated against
  Anki's own importer).
- **Hover-popup word lookup** (#780): hovering a word in a definition for ~0.4s
  pops up its full definition near the cursor (macOS Dictionary / GoldenDict
  style). Double-click lookup (#258) and `entry://` jumps are unchanged. Toggle
  via the `hoverpopup` preference.
- **Export current entry as HTML** (#784): a toolbar button saves the currently
  displayed entry as a self-contained HTML file (resources inlined) for debugging
  complex entries.
- **Settings write-back** (#777): user preferences now persist to `medict.toml`
  (the foundation used by multi-dict selection and the hover-popup toggle).

### Changed
- **Bookmarks v2** (#772): the 生词本 store migrated to a pure-Go SQLite store
  with notebooks (create / rename / delete / set-default) and **offline HTML
  snapshots** — saved words stay viewable even after the source dictionary is
  unloaded.

### Fixed
- **Dictionary display name** (#782, #788): the dict list now shows the real
  dictionary title (mdx header `Title` / stardict `bookname`) instead of the
  file name — including at first load, before the index is built.

## v3.1.6

Header UI redesign + bookmarks nav-rail fix.

### Fixed
- **Bookmarks page missing navigation** (#768): 生词本 page lacked `<AppHeader>`,
  making it impossible to switch tabs. Fixed.
- **macOS double title bar**: removed the redundant 26px fake-title-bar
  (`Frameless:false` already provides a native title bar).

### Changed
- **Header redesign**: logo shrunk (h1 24px→16px, area 120px→72px), search bar
  now `flex:1` (fills width), back/forward/star buttons restyled (transparent bg,
  hover, gap spacing), header height 60px→48px.

## v3.1.5

Hotfix: bookmarks page navigation.

### Fixed
- **Bookmarks page missing navigation rail** (#768): the 生词本 page was missing
  `<AppHeader>` (which contains `AppFunctions`), so users could not switch tabs
  after entering it. Now wrapped in the standard `x-layout` shell.

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
