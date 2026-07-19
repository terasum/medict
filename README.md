# Medict

[![Latest Release](https://img.shields.io/github/v/release/terasum/medict?display_name=tag&style=flat-square)](https://github.com/terasum/medict/releases/latest)
[![Release](https://github.com/terasum/medict/actions/workflows/package-and-release.yaml/badge.svg)](https://github.com/terasum/medict/actions/workflows/package-and-release.yaml)
[![CI](https://github.com/terasum/medict/actions/workflows/ci.yml/badge.svg)](https://github.com/terasum/medict/actions/workflows/ci.yml)
[![GitHub Stars](https://img.shields.io/github/stars/terasum/medict?style=flat-square)](https://github.com/terasum/medict/stargazers)
[![License](https://img.shields.io/github/license/terasum/medict.svg?style=flat-square)](LICENSE)

Medict 是一款开源、跨平台的本地词典应用，支持 MDict（`.mdx` / `.mdd`）和 StarDict 词典。它基于 Wails、Go 与 Vue 构建，提供离线英汉词典、多词典查询、生词本、Anki 导出和词典样式编辑等完整的桌面查词体验。

当前稳定版本为 **v3.1.11**。

## 主要功能

### 查词与阅读

- 支持 MDict v1.x / v2.0 的 `.mdx`、`.mdd` 词典，以及 StarDict 的 `.ifo`、`.idx`、`.dict` / `.dict.dz` 词典。
- 内置约 5 万高频词的离线 ECDICT 英汉词典，无需账号、密钥或网络即可使用。
- 可在「设置 → 词典设置」中下载完整 ECDICT 数据，安装后继续完全离线查询。
- 支持前缀建议与模糊搜索，拼写不完全准确时仍可找到候选词。
- 支持单词典查询和多词典同时查询；多词典选择会自动保存。
- 正确加载词典内的 CSS、JavaScript、图片、字体与音频资源。
- 支持 `entry://` 词条跳转和 `@@@LINK=` 重定向。
- 鼠标悬停时自动识别并为单词添加下划线，点击后在主页面查询，体验接近 macOS 词典。
- 支持前进、后退、刷新，以及工具栏、`Cmd/Ctrl + +/-`、`Cmd/Ctrl + 滚轮`缩放释义内容。
- 缩放范围为 80%–160%，会应用到单词典和多词典内容并自动保存。

### 词典管理

- 自动递归扫描词典目录，并支持通过软链接引入词典。
- 自动读取 MDict `Title` 或 StarDict `bookname` 作为显示名称。
- 优先显示词典自带封面；没有封面时使用统一的默认图标。
- 可创建词典组、管理组成员、删除分组和标记常用分组，不会改动原始词典文件。
- 支持在独立窗口中编辑单个词典的 CSS；已有样式会作为编辑基础，并支持实时预览和持久化覆盖。
- 可将当前词条导出为资源内联的独立 HTML，便于调试复杂词典内容。

### 生词本

- 查询时点击星标即可收藏当前词条。
- 支持创建、重命名、删除生词本和设置默认生词本。
- 收藏时保存离线 HTML 快照，即使原词典之后被移除，仍可查看释义。
- 可按词条或来源词典筛选，并从生词列表直接返回查询。
- 可将整个生词本导出为 Anki `.apkg`，词典图片会一并打包为 Anki 媒体资源。

### 设置与调试

- 统一管理词典目录、完整 ECDICT、软件、主题和插件相关设置。
- 使用说明、隐私声明和开源许可直接集成在设置页面中。
- 内置词典资源查询调试工具，可按词典检查图片、CSS、字体或音频资源。
- 「关于软件」集中提供版本信息、源代码、版本发布和问题反馈入口。

## 下载

请从 [GitHub Releases](https://github.com/terasum/medict/releases/latest) 下载最新稳定版本。每个正式版本由 GitHub Actions 自动构建并提供以下产物：

| 平台 | 构建产物 |
| --- | --- |
| macOS（Apple Silicon / Intel） | Universal `.dmg`、`.app.zip` |
| Windows x86_64 | `.zip` |
| Linux x86_64 | `.tar.gz` |

历史版本和完整更新记录请参阅 [CHANGELOG.md](CHANGELOG.md)。

## 安装词典

Medict 以目录为单位识别词典。打开「设置 → 词典设置」，点击词典目录旁的按钮即可进入当前系统实际使用的目录。将每部词典放入独立子目录后，重新启动或重新加载应用即可。

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

目录可以继续按语言或用途嵌套，Medict 会递归找到其中真正包含词典文件的目录。

### 词典文件规则

| 类型 | 必需文件 | 可选文件 |
| --- | --- | --- |
| MDict | 一个 `.mdx` | 同目录下一个或多个 `.mdd`；与 `.mdx` 同名的 `.jpg` / `.png` 封面 |
| StarDict | `.ifo`、`.idx`、`.dict` 或 `.dict.dz` | `cover.jpg` / `cover.png` |

通用封面也可以命名为 `cover.jpg` 或 `cover.png`。词典内容及其授权由词典提供者负责，请仅使用你有权使用的词典文件。

## 从源码运行

### 环境要求

- Go 1.25
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)
- [Bun](https://bun.sh/)
- macOS、Windows，或安装了 GTK3 / WebKitGTK 的 Linux 环境

安装 Wails CLI：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

安装前端依赖并启动开发模式：

```bash
cd frontend
bun install --frozen-lockfile
cd ..
make dev
```

构建桌面应用：

```bash
make build
```

产物会生成在 `build/bin/`。

## 测试

```bash
# Go 全量测试
go test ./...

# 前端测试、类型检查和生产构建
cd frontend
bun run test
bun run typecheck
bun run build
```

## 技术栈

- 桌面框架：[Wails v2](https://wails.io/)
- 后端：Go、Gin、LevelDB、SQLite
- 前端：Vue 3、Vite 6、Pinia、Naive UI、CodeMirror 6
- 测试：Go testing、Vitest、Vue Test Utils

## 参与贡献

欢迎通过 [Issues](https://github.com/terasum/medict/issues) 报告问题、提交词典兼容性案例或提出功能建议。提交界面问题时，请尽量附上操作系统、Medict 版本、复现步骤和截图；提交词典解析问题时，请说明词典格式与可公开的最小复现信息。

## 许可证

Medict 使用 [GNU General Public License v3.0](LICENSE) 发布。

**Medict is made by terasum and xing with ❤️**
