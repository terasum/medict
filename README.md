# Medict
![Static Badge](https://img.shields.io/badge/version-v3.0.1-blue?style=flat)
![GitHub Repo stars](https://img.shields.io/github/stars/terasum/medict)
![license](https://img.shields.io/github/license/terasum/medict.svg)

## Medict 介绍

Medict 是一个跨平台的词典 APP, 主要支持 \*.mdx/\*.mdd 词典格式, 目前支持 v1.x 和 v2.0 格式的词典。

## v3-ing
Medict v3 版本正在开发中，采用 wails 框架重构，启动性能与包容量大小将大大优化，尽情期待。
Medict version 3 is under developing, will refactor by wails framework, waiting for time!

## V3-index
<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/v3-medict-app-index.png" alt="v3词典界面" style="zoom: 23%;" />
</div>


## 下载与更新

目前 Medict 正在紧张开发阶段，版本为自动打包滚动发布，请自行到 https://github.com/terasum/medict/releases 页面寻找最新开发版本, 所有版本均有打包日期，选择最新版本即可。

## 使用说明

### 界面概览

软件目前包括 “搜索”，“词典”，“插件”，“设置” 四个界面，其中“插件”目前尚在开发当中。
<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/zov3hq.png" alt="v3词典界面" style="zoom: 23%;" />
</div>

### 词典安装

目前v3版本的Medict采用自动扫描词典的方式添加词典，您只需将词典放到词典扫描根目录即可，以下是各个操作系统的词典默认扫描根目录：

| 操作系统 | 默认词典目录                                              | 说明                                                         |
| -------- | --------------------------------------------------------- | ------------------------------------------------------------ |
| MacOS    | $HOME/Library/Application\ Support/medict/dicts | 您可以在Medict 的`设置`界面打开，也可打开Finder,按住 `Control+Shift+G`并输入该目录打开。 |
| Windows  | %APPDATA%/medict/dicts                                    |                                                              |
| Linux    | $HOME/.medict/dicts                                       |                                                              |

目前 Medict 采用的是目录词典组织方式，一个目录即为一个词典，内部可以放置mdict或stardict两种格式的词典，类似下图：

<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/image-20230927112451745.png" alt="v3词典界面" style="zoom: 23%;" />
</div>

#### 特殊文件

在目录中可以放置一些特殊用途的文件，用于增强词典：

| 特殊文件名称               | 用途         | 说明                                                         |
| -------------------------- | ------------ | ------------------------------------------------------------ |
| cover.jpg或cover.png       | 词典封面图片 | 该文件针对stardict词典，若是mdict词典，也可提供与mdx同名的jpg或png文件 |
| mdict.dtype/stardict.dtype | 词典类型文件 | 后续用于支持其他词典类型（目前可自动识别）                   |

#### stardict

目前仅支持同时具备 (dict.dz或dict文件，ifo文件，idx文件)的词典，缺少文件的词典不支持。

### 查词功能

在搜索框输入所需的词即可查词，右侧边栏可以选择使用的词典：


<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/image-20230927112656769.png" alt="v3词典界面" style="zoom: 23%;" />
</div>


### 词典界面

点击上方词典按钮，可进入词典功能，进行词典管理，目前仅支持查看，词典请通过操作系统文件浏览器管理（Finder等）


<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/image-20230927113137569.png" alt="v3词典界面" style="zoom: 23%;" />
</div>

### 插件管理

**tobe done**



### 词典设置

点击上方设置按钮，可进入词典设置功能。

<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="docs/_assets/image-20230927113331495-5785613-5785617.png" alt="v3词典界面" style="zoom: 23%;" />
</div>

## Q&A

### 发音问题

目前 oale8 词典这种内嵌发音按钮的，将音频资源嵌入在mdd文件中的词典是可以支持发音的，但是目前采用的是js替换的方式完成，不一定适用于所有词典，需要case by case 调试


## Call for help
**目前 Medict 正在开发当中，需要您的帮助！**

目前希望得到的帮助：
1. UI 设计 / Logo 设计
2. e2e 测试框架集成
3. 词典测试
4. 词典内容安全测试

## 特性列表
- [x] APP 新 LOGO
- [x] 查词建议 suggest list
- [x] mdx 查词结果展示
- [x] mdd 资源加载
- [x] mdd 音频播放(mp3/ogg)
- [x] entry:// 词汇跳转
- [x] @@@LINK= 词汇重定向
  - [ ] 存在部分词典适配问题
- [x] mdd/mdx 词典选择配置
  - [x] mdd 可选配置
- [ ] 查词历史导航(</>)
- [ ] 功能 tab 页跳转
- [ ] 多词典同时查询
- [ ] 全文检索^1
- [ ] 有道等在线词库增强
- [ ] 插件功能
  - [ ] 词典扩展功能栏
  - [ ] 词频展示插件
  - [ ] 生词本记录插件
  - [ ] 导出anki卡片插件

**Medict is made by terasum and xing with ❤️**
