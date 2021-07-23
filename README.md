# Medict
[![NodeCI](https://github.com/terasum/medict/workflows/Node%20CI/badge.svg?event=push)](https://github.com/terasum/medict/actions?query=workflow%3A%22Node+CI%22+branch%3Acanary+event%3Apush)
![GitHub release](https://img.shields.io/github/package-json/v/terasum/medict)
![license](https://img.shields.io/github/license/terasum/medict.svg)

Medict 是一个跨平台的词典 APP, 主要支持 \*.mdx/\*.mdd 词典格式, 目前支持 v1.x和v2.0格式的词典

**目前 Medict 正在开发当中，需要您的帮助！**

目前希望得到的帮助：
1. UI 设计 / Logo 设计
2. e2e 测试框架集成
3. 词典测试
4. 词典内容安全测试

## 特性列表
- [ ] APP 新 LOGO
- [x] 查词建议 suggest list
- [x] mdx 查词结果展示
- [x] mdd 资源加载
- [x] mdd 音频播放(mp3/ogg)
- [x] entry:// 词汇跳转
- [x] @@@LINK= 词汇重定向
- [ ] mdd/mdx 词典选择配置
- [ ] 查词历史导航(</>)
- [ ] 功能 tab 页跳转
- [ ] 多词典同时查询
- [ ] 全文检索^1
- [ ] 有道等在线词库增强
- [ ] 翻译功能^2
- [ ] 插件功能
  - [ ] 词典扩展功能栏
  - [ ] 词频展示插件
  - [ ] 生词本记录插件
  - [ ] 导出anki卡片插件


## UI

界面预览

![preview](docs/images/medict-capture.gif)


## 使用方法

目前版本仅提供开发预览版本，目前尚未编译安装版，仅供开发人员阅读

### 克隆代码

``` bash
git clone https://github.com/terasum/medict.git
```

### 开发模式运行

```
cd medict
yarn install
yarn start
```

**Medict is made by terasum with ❤️**
