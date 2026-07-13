# Changelog

All notable changes to [Medict](https://github.com/terasum/medict) are recorded
here. The most recent release is at the top. For the full history before
v3.1.0, see `git log v3.0.1..HEAD`.

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
