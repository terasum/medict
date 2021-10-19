# Medict
[![NodeCI](https://github.com/terasum/medict/workflows/Node%20CI/badge.svg?event=push)](https://github.com/terasum/medict/actions?query=workflow%3A%22Node+CI%22+branch%3Acanary+event%3Apush)
![GitHub release](https://img.shields.io/github/package-json/v/terasum/medict)
![license](https://img.shields.io/github/license/terasum/medict.svg)

## Medict 介绍

Medict 是一个跨平台的词典 APP, 主要支持 \*.mdx/\*.mdd 词典格式, 目前支持 v1.x 和 v2.0 格式的词典。

## 下载与更新

目前 Medict 正在紧张开发阶段，版本为自动打包滚动发布，请自行到 https://github.com/terasum/medict/releases 页面寻找最新开发版本, 所有版本均有打包日期，选择最新版本即可。

## 界面概览

软件目前包括 “词典”，“翻译”，“插件”，“设置” 四个界面，其中“插件”目前尚在开发当中。

### 词典界面

<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="src/renderer/assets/docs/pic_dict_window.jpg" alt="词典界面" style="zoom: 23%;" />
</div>

### 翻译界面

<div style="width: 100%;">
<img  width=500 style="display:block; margin: 0 auto;"  src="src/renderer/assets/docs/pic_translate_window.jpg" alt="翻译界面" style="zoom: 23%;" />
</div>

### 设置界面

<div style="width: 100%;">
  <img  width=500 style="display:block; margin: 0 auto;"  src="src/renderer/assets/docs/pic_settings_window.jpg" alt="设置界面" style="zoom: 23%;" />
</div>

## 下载

目前Medict正在紧张开发阶段，版本为自动打包滚动发布，请自行到 [release](https://github.com/terasum/medict/releases) 页面寻找最新开发版本, 所有版本均有打包日期，选择最新版本即可。


## 使用步骤
### 步骤1: 添加词典

1. 点击右上角设置
2. 点击下方 "+" 号

<div style="width: 100%;">
    <img width=500 style="display:block; margin: 0 auto;" src="src/renderer/assets/docs/pic_add_dict_btn.jpg" alt="pic_add_dict_btn.jpg" style="zoom: 23%;" />
</div>

3. 在弹出框中填写词典信息
4. 选择词典文件

<div>
  <img width=500 style="display:block; margin: 0 auto;"  src="src/renderer/assets/docs/pic_add_dict_modal.jpg" alt="pic_add_dict_modal.jpg" style="zoom:23%;" />
</div>

**注意:** mdx文件所在的文件夹中的js/css/font文件均会被拷贝到缓存文件夹中，请把一个独立词典放在一个独立的文件夹中，并将相关资源放在一起。

**注意:** mdx/mdd 本身不会被拷贝，删除之后，词典将无法找到该mdx文件

## 步骤 2: 查词

1. 选择词典并输入目标词（模糊）

<div>
  <img width=500 style="display:block; margin: 0 auto;" src="src/renderer/assets/docs/pic_usage_step1.jpg" alt="pic_usage_step1.jpg" style="zoom:23%;" />
</div>

2. 在左边栏选择你想要查的具体词汇, 如果该词汇和其他词汇同一个意思（即@@Link==） 则直接展示该同意义词汇

<div>
<img width=500 style="display:block; margin: 0 auto;" src="src/renderer/assets/docs/pic_usage_step2.jpg" alt="pic_usage_step2.jpg" style="zoom:23%;" />
</div>

## Q&A

### 发音问题

目前 oale8 词典这种内嵌发音按钮的，将音频资源嵌入在mdd文件中的词典是可以支持发音的，但是目前采用的是js替换的方式完成，不一定适用于所有词典，需要case by case 调试


### 跳转问题

目前有两种跳转：

1. @@Link 的跳转，自动跳转，但是如果出现跳转环路，会停止跳转，直接展示 @@Link==
2. 内部跳转即 `<a href="entry://">` 的方式，如果entry中间是完整词条，可以支持跳转，如果是特殊词条，目前还不支持

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
- [x] 翻译功能^2
- [ ] 插件功能
  - [ ] 词典扩展功能栏
  - [ ] 词频展示插件
  - [ ] 生词本记录插件
  - [ ] 导出anki卡片插件


## 开发说明

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

**Medict is made by terasum and xing with ❤️**
