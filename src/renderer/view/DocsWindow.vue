<template>
  <div class="container-fluid" style="height: 100%; width: 100%">
    <Header :displaySearchBox="true" />
    <div class="window-content">
      <div class="pane pane-sm sidebar">
        <div class="nav-group">
          <h5 class="nav-group-title">基本使用说明</h5>
          <span
            class="nav-group-item"
            :class="currentMenu === 0 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(0)"
          >
            主要界面介绍
          </span>
          <span
            class="nav-group-item"
            :class="currentMenu === 1 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(1)"
          >
            词典配置使用说明
          </span>
          <span
            class="nav-group-item"
            :class="currentMenu === 2 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(2)"
          >
            百度翻译配置
          </span>

          <span
            class="nav-group-item"
            :class="currentMenu === 3 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(3)"
          >
            有道翻译配置
          </span>
          <span
            class="nav-group-item"
            :class="currentMenu === 4 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(4)"
          >
            FAQ 
          </span>

          <h5 class="nav-group-title">软件协议与信息</h5>
          <span
            class="nav-group-item"
            :class="currentMenu === 5 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(5)"
          >
            免责声明
          </span>

          <span
            class="nav-group-item"
            :class="currentMenu === 6 ? 'item-active' : ''"
            @click="onClickPreferenceMenu(6)"
          >
            开源协议信息
          </span>
        </div>
      </div>
      <div class="docs-container">
        <div class="markdown-body" v-html="docs"></div>
      </div>
    </div>
    <FooterBar />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';
// @ts-ignore
import index_md from '../assets/docs/index.md';
// @ts-ignore
import select_and_use_md from '../assets/docs/select_and_use_dict.md';
// @ts-ignore
import baidu_translate_md from '../assets/docs/baidu_translate_config.md';
// @ts-ignore
import youdao_translate_md from '../assets/docs/youdao_traslate_config.md';
// @ts-ignore
import faq_md from '../assets/docs/faq.md';
// @ts-ignore
import terms_and_service from '../assets/docs/terms_and_service.md';
// @ts-ignore
import license_md from '../assets/docs/license.md';

const routerMap = {
  0: index_md,
  1: select_and_use_md,
  2: baidu_translate_md,
  3: youdao_translate_md,
  4: faq_md,
  5: terms_and_service,
  6: license_md,
};

export default Vue.extend({
  components: { Header, FooterBar },
  data() {
    return {
      docs: '',
      currentMenu: 0,
    };
  },
  methods: {
    onClickPreferenceMenu(id: number) {
      console.log(`click id ${id} router: ${routerMap[id]}`);
      if (!routerMap[id]) {
        return;
      }
      if (this.currentMenu == id) {
        return;
      }
      this.currentMenu = id;
      this.docs = routerMap[id];
    },
  },
  mounted() {
    this.$nextTick(() => {
      this.docs = index_md;
    });
  },
});
</script>

<style lang="scss" scoped>
.window-content {
  height: calc(100% - 80px);
  width: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: row;
  .pane-group {
    width: 100%;
    display: flex;
  }
  .sidebar {
    width: 160px;
    height: 100%;
    border-right: 1px solid #e8e8e8;
    background: #f2f4f5;

    .nav-group {
      user-select: none;
      display: flex;
      flex-direction: column;
      padding: 10px 6px;

      .nav-group-title {
        width: 100%;
        margin: 0;
        padding-left: 6px;
        font-size: 12px;
        font-weight: 500;
        color: #999;
        border-bottom: 1px solid #c1c1c1;
        margin-bottom: 4px;
        margin-top: 6px;
        &:nth-child(1) {
          margin-top: 0px;
        }
      }

      .nav-group-item {
        height: 26px;
        width: 100%;
        display: flex;
        color: #777;
        text-decoration: none;
        font-size: 12px;
        line-height: 26px;
        cursor: pointer;
        padding-left: 20px;

        &:active {
          background-color: #e8eaec;
        }

        .icon {
          height: 30px;
          width: 30px;
          display: inline-block;
          font-size: 14px;
          line-height: 30px;
          text-align: center;
          color: #777;
        }
      }
      .item-active {
        border-radius: 3px;
        background-color: #e8eaec;
        color: #222;
      }
    }
  }

  .docs-container {
    height: 100%;
    margin: 0;
    overflow-y: auto;
    width: calc(100% - 160px);
    padding: 4px 10px;

    .markdown-body {
      font-size: 100%;
      overflow-y: scroll;
      -webkit-text-size-adjust: 100%;
      -ms-text-size-adjust: 100%;
      background: #fefefe;

      color: #444;
      font-family: Georgia, Palatino, 'Palatino Linotype', Times,
        'Times New Roman', serif;
      font-size: 14px;
      line-height: 1.5em;
      padding: 1em;
      margin: auto;
      max-width: 42em;
      background: #fefefe;
    }
  }
}
</style>

<style lang="scss">
.markdown-body {
  a {
    color: #0645ad;
    text-decoration: none;
  }
  a:visited {
    color: #0b0080;
  }
  a:hover {
    color: #06e;
  }
  a:active {
    color: #faa700;
  }
  a:focus {
    outline: thin dotted;
  }
  a:hover,
  a:active {
    outline: 0;
  }

  ::-moz-selection {
    background: rgba(255, 255, 0, 0.3);
    color: #000;
  }
  ::selection {
    background: rgba(255, 255, 0, 0.3);
    color: #000;
  }

  a::-moz-selection {
    background: rgba(255, 255, 0, 0.3);
    color: #0645ad;
  }
  a::selection {
    background: rgba(255, 255, 0, 0.3);
    color: #0645ad;
  }

  p {
    margin: 1em 0;
  }

  img {
    max-width: 100%;
    margin: 1em 0;
  }

  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    font-weight: normal;
    color: #343434;
    line-height: 1em;
    margin: 1em 0;
  }
  h4,
  h5,
  h6 {
    font-weight: bold;
  }
  h1 {
    font-size: 1.8em;
    text-align: center;
  }
  h2 {
    font-size: 1.5em;
  }
  h3 {
    font-size: 1.3em;
  }
  h4 {
    font-size: 1.1em;
  }
  h5 {
    font-size: 1em;
  }
  h6 {
    font-size: 1em;
  }

  blockquote {
    color: #666666;
    margin: 0;
    padding-left: 3em;
    border-left: 0.5em #eee solid;
  }
  hr {
    display: block;
    height: 0;
    border: 0;
    border-top: 1px solid #aaa;
    border-bottom: 1px solid #eee;
    margin: 1em 0;
    padding: 0;
  }
  pre,
  code,
  kbd,
  samp {
    color: #000;
    font-family: monospace, monospace;
    _font-family: 'courier new', monospace;
    font-size: 0.98em;
  }
  pre {
    white-space: pre;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  b,
  strong {
    font-weight: bold;
  }

  dfn {
    font-style: italic;
  }

  ins {
    background: #ff9;
    color: #000;
    text-decoration: none;
  }

  mark {
    background: #ff0;
    color: #000;
    font-style: italic;
    font-weight: bold;
  }

  sub,
  sup {
    font-size: 75%;
    line-height: 0;
    position: relative;
    vertical-align: baseline;
  }
  sup {
    top: -0.5em;
  }
  sub {
    bottom: -0.25em;
  }

  ul,
  ol {
    margin: 1em 0;
    padding: 0 0 0 2em;
  }
  li p:last-child {
    margin: 0;
  }
  dd {
    margin: 0 0 0 2em;
  }

  img {
    border: 0;
    -ms-interpolation-mode: bicubic;
    vertical-align: middle;
  }

  table {
    border-collapse: collapse;
    border-spacing: 0;
  }
  td {
    vertical-align: top;
  }
}
</style>