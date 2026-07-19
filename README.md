# Medict

**English** | [简体中文](README_ZH-CN.md)

[![Latest Release](https://img.shields.io/github/v/release/terasum/medict?display_name=tag&style=flat-square)](https://github.com/terasum/medict/releases/latest)
[![Release](https://github.com/terasum/medict/actions/workflows/package-and-release.yaml/badge.svg)](https://github.com/terasum/medict/actions/workflows/package-and-release.yaml)
[![CI](https://github.com/terasum/medict/actions/workflows/ci.yml/badge.svg)](https://github.com/terasum/medict/actions/workflows/ci.yml)
[![GitHub Stars](https://img.shields.io/github/stars/terasum/medict?style=flat-square)](https://github.com/terasum/medict/stargazers)
[![License](https://img.shields.io/github/license/terasum/medict.svg?style=flat-square)](LICENSE)

Medict is an open-source, cross-platform desktop dictionary application with support for MDict (`.mdx` / `.mdd`) and StarDict dictionaries. Built with Wails, Go, and Vue, it provides an offline English-Chinese dictionary, multi-dictionary lookup, vocabulary notebooks, Anki export, and per-dictionary style editing.

The current stable release is **v3.1.11**.

## Features

### Lookup and Reading

- Supports MDict v1.x / v2.0 `.mdx` and `.mdd` dictionaries, plus StarDict `.ifo`, `.idx`, `.dict`, and `.dict.dz` dictionaries.
- Includes an offline ECDICT English-Chinese dictionary with about 50,000 high-frequency entries—no account, API key, or network connection required.
- The complete ECDICT dataset can be downloaded from **Settings → Dictionary Settings** and remains fully available offline after installation.
- Provides prefix suggestions and fuzzy search, helping users find candidates even when spelling is not exact.
- Supports both single-dictionary and simultaneous multi-dictionary lookup; the active selection is persisted automatically.
- Loads dictionary-provided CSS, JavaScript, images, fonts, and audio resources.
- Supports `entry://` navigation and `@@@LINK=` redirects.
- Detects words under the pointer, underlines them on hover, and looks them up in the main view when clicked, similar to macOS Dictionary.
- Supports back, forward, and refresh actions, plus definition zoom through the toolbar, `Cmd/Ctrl + +/-`, or `Cmd/Ctrl + mouse wheel`.
- Definition zoom ranges from 80% to 160%, applies to single- and multi-dictionary views, and is persisted automatically.

### Dictionary Management

- Recursively scans the dictionary directory and supports dictionaries referenced through symbolic links.
- Uses the MDict `Title` or StarDict `bookname` metadata as the display name when available.
- Prefers dictionary-provided cover artwork and falls back to a shared default icon.
- Supports creating dictionary groups, managing group members, deleting groups, and marking a favorite group without modifying dictionary files.
- Provides a standalone CSS editor for each dictionary. Existing dictionary styles are used as the editing base, with live preview and persistent overrides.
- Exports the current entry as a self-contained HTML file with resources inlined for debugging complex dictionary content.

### Vocabulary Notebooks

- Save the current entry by clicking the star while looking up a word.
- Create, rename, delete, and select a default vocabulary notebook.
- Stores an offline HTML snapshot when an entry is saved, so definitions remain available even if the source dictionary is later removed.
- Filter saved words by entry or source dictionary and return directly to lookup from the notebook.
- Export an entire notebook as an Anki `.apkg`; dictionary images are packaged as Anki media resources.

### Settings and Diagnostics

- Manage the dictionary directory, complete ECDICT dataset, software, theme, and plugin-related preferences in one place.
- Read the user guide, privacy notice, and open-source license directly inside Settings.
- Inspect images, CSS, fonts, and audio through the built-in dictionary resource query tool.
- Access version information, source code, releases, and issue reporting from the About page.

## Download

Download the latest stable build from [GitHub Releases](https://github.com/terasum/medict/releases/latest). Every official release is built automatically by GitHub Actions and includes:

| Platform | Artifact |
| --- | --- |
| macOS (Apple Silicon / Intel) | Universal `.dmg` and `.app.zip` |
| Windows x86_64 | `.zip` |
| Linux x86_64 | `.tar.gz` |

See [CHANGELOG.md](CHANGELOG.md) for the complete release history.

## What's New in v3.1.11

### Added

- Dictionary groups in the dictionary sidebar, including group creation, member management, deletion, and favorite-group selection.
- Persistent definition zoom from 80% to 160% in both single- and multi-dictionary modes.
- A standalone dictionary CSS editor that starts from existing styles when available and retains live preview.
- Top-level documentation tabs in Settings for the user guide, privacy notice, and open-source license.

### Improved

- Switched the default English-Chinese experience to offline-first: the bundled ECDICT subset remains available, with an in-app installer for the complete dataset.
- Reorganized Settings and diagnostics into a flatter structure, with resource queries integrated directly into Settings.
- Replaced popup-based definition lookup with hover underlining followed by click-to-lookup.
- Unified top navigation, dictionary icons and active states, sidebars, bottom action bars, settings controls, and footer styling.
- Upgraded the frontend build toolchain to Vite 6.

### Fixed

- Definition jumps and sidebar lookups now keep the search field synchronized, and entries remain selectable after navigation.
- Dictionary-provided cover artwork is preferred, with a shared fallback icon when no artwork exists.
- Online dictionary errors now degrade gracefully instead of returning unhandled HTTP 500 responses.

See the [v3.1.11 release](https://github.com/terasum/medict/releases/tag/v3.1.11) or read the [complete changelog](CHANGELOG.md).

## Roadmap

The following areas are currently planned. Priorities may change based on user feedback and contributions:

- [ ] [Full-text search](https://github.com/terasum/medict/issues/787): search inside all definitions instead of matching headwords only.
- [ ] [Inflection, Simplified/Traditional Chinese, and variant-character fallback](https://github.com/terasum/medict/issues/786): retry a missing lookup using English inflections and Chinese character variants.
- [ ] [Performance and stability for large dictionaries](https://github.com/terasum/medict/issues/781): continue investigating indexing, lookup, and rendering stalls in some dictionaries.
- [ ] [MDX packing and unpacking tools](https://github.com/terasum/medict/issues/785): provide a local toolchain for dictionary maintainers.
- [ ] [Custom resource extraction directory](https://github.com/terasum/medict/issues/120): allow users to select where dictionary resources are cached or extracted.

See [GitHub Issues](https://github.com/terasum/medict/issues) for more plans, bugs, and discussions.

## Installing Dictionaries

Medict treats each directory as one dictionary. Open **Settings → Dictionary Settings**, then use the button next to the dictionary directory to open the actual location used by your system. Put each dictionary in its own subdirectory, then restart or reload the application.

```text
dicts/
├── Oxford/
│   ├── Oxford.mdx
│   ├── Oxford.mdd
│   └── Oxford.jpg
└── StardictExample/
    ├── example.ifo
    ├── example.idx
    └── example.dict.dz
```

Directories may be nested by language or purpose. Medict recursively discovers the directories that actually contain dictionary files.

### Dictionary File Rules

| Type | Required files | Optional files |
| --- | --- | --- |
| MDict | One `.mdx` | One or more `.mdd` files in the same directory; `.jpg` / `.png` cover with the same base name as the `.mdx` |
| StarDict | `.ifo`, `.idx`, and either `.dict` or `.dict.dz` | `cover.jpg` / `cover.png` |

A shared cover may also be named `cover.jpg` or `cover.png`. Dictionary content and licensing remain the responsibility of each dictionary provider; only use dictionary files you are authorized to use.

## Running from Source

### Requirements

- Go 1.25
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)
- [Bun](https://bun.sh/)
- macOS, Windows, or Linux with GTK3 and WebKitGTK installed

Install the Wails CLI:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Install frontend dependencies and start development mode:

```bash
cd frontend
bun install --frozen-lockfile
cd ..
make dev
```

Build the desktop application:

```bash
make build
```

Build artifacts are written to `build/bin/`.

## Testing

```bash
# Run all Go tests
go test ./...

# Run frontend tests, type checking, and the production build
cd frontend
bun run test
bun run typecheck
bun run build
```

## Technology Stack

- Desktop framework: [Wails v2](https://wails.io/)
- Backend: Go, Gin, LevelDB, SQLite
- Frontend: Vue 3, Vite 6, Pinia, Naive UI, CodeMirror 6
- Testing: Go testing, Vitest, Vue Test Utils

## Contributing

Bug reports, dictionary compatibility cases, and feature proposals are welcome through [GitHub Issues](https://github.com/terasum/medict/issues). For interface issues, please include your operating system, Medict version, reproduction steps, and screenshots when possible. For dictionary parsing issues, describe the dictionary format and provide the smallest publicly shareable reproduction details.

## License

Medict is released under the [GNU General Public License v3.0](LICENSE).

**Medict is made by terasum and xing with ❤️**
