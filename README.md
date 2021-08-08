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


## 界面预览



### 英汉

**OALE8**

- 支持图片展开
- 支持 jquery
- 支持跳转

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst1zj6bdbj313n0u0grs.jpg" alt="image-20210725115204014" style="zoom: 33%;" />



### 汉英

**新世纪汉英大辞典**

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst22mzl8cj313n0u0n3c.jpg" alt="image-20210725115503487" style="zoom: 33%;" />

### 日语

**新日汉大辞典**

- 支持字体

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst1vdd6lyj613n0u0tef02.jpg" alt="image-20210725114802241" style="zoom: 33%;" />



## 图片词库

**大英漢词典**

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst1wfn7ffj313n0u0qd3.jpg" alt="image-20210725114905767" style="zoom: 33%;" />



## 下载

目前Medict正在紧张开发阶段，版本为自动打包滚动发布，请自行到 [release](https://github.com/terasum/medict/releases) 页面寻找最新开发版本, 所有版本均有打包日期，选择最新版本即可。


## 使用步骤

### 步骤1: 添加词典

1. 点击右上角设置
2. 点击下方 "+" 号

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst2568g8dj313n0u0wi1.jpg" alt="image-20210725115729944" style="zoom: 33%;" />



3. 在弹出框中填写词典信息
4. 选择词典文件

注意：mdx文件所在的文件夹中的js/css/font文件均会被拷贝到缓存文件夹中，请把一个独立词典放在一个独立的文件夹中，并将相关资源放在一起。

注意： mdx/mdd 本身不会被拷贝，删除之后，词典将无法找到该mdx文件



<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst27etnrgj613n0u0afx02.jpg" alt="image-20210725115939188" style="zoom:33%;" />



## 步骤2: 查词

1. 选择词典并输入目标词（模糊）

<img src="https://tva1.sinaimg.cn/large/008i3skNly1gst2anwbx3j313n0u0jz4.jpg" alt="image-20210725120246204" style="zoom:33%;" />



2. 在左边栏选择你想要查的具体词汇

   如果该词汇和其他词汇同一个意思（即@@Link==） 则直接展示该同意义词汇

   <img src="https://tva1.sinaimg.cn/large/008i3skNly1gst2katvp6j313n0u010a.jpg" alt="image-20210725121201887" style="zoom:33%;" />



## Q&A

### 发音问题

目前 oale8 词典这种内嵌发音按钮的，将音频资源嵌入在mdd文件中的词典是可以支持发音的，但是目前采用的是js替换的方式完成，不一定适用于所有词典，需要case by case 调试



### 跳转问题

目前有两种跳转：

1. @@Link 的跳转，自动跳转，但是如果出现跳转环路，会停止跳转，直接展示 @@Link==
2. 内部跳转即 `<a href="entry://">` 的方式，如果entry中间是完整词条，可以支持跳转，如果是特殊词条，目前还不支持


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
