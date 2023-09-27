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

注意!!!
v3 版本正在开发中，请切换分支至 master 再进行下载老版本。


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
- [ ] entry:// 词汇跳转
- [ ] @@@LINK= 词汇重定向
  - [ ] 存在部分词典适配问题
- [ ] mdd/mdx 词典选择配置
  - [ ] mdd 可选配置
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
