<template>
  <div class="container" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="translate-container" style="height: 100%">
      <div class="sidebar">
        <!-- translate engines selection -->
        <div class="translator-engine">
          <b-dropdown aria-role="list" id="translator-engine">
            <template #trigger="{}">
              <b-button
                ><span class="lang-icon"><i class="fas fa-language"></i></span>
                {{ selectedEngine }}</b-button
              >
            </template>
            <b-dropdown-item @click="useEngine('google')">谷歌</b-dropdown-item>
            <b-dropdown-item @click="useEngine('baidu')">百度</b-dropdown-item>
            <b-dropdown-item @click="useEngine('youdao')">有道</b-dropdown-item>
          </b-dropdown>
        </div>

        <!-- source language selection -->
        <div class="translator-src-lang">
          <b-dropdown aria-role="list" id="src-lang">
            <template #trigger="{}">
              <b-button :label="sourceLang" type="is-primary" />
            </template>

            <b-dropdown-item @click="changeSourceLang('en')"
              >英文</b-dropdown-item
            >
            <b-dropdown-item @click="changeSourceLang('zh')"
              >中文</b-dropdown-item
            >
            <b-dropdown-item @click="changeSourceLang('jp')"
              >日文</b-dropdown-item
            >
          </b-dropdown>
        </div>

        <div class="translator-icon-container">
          <span>
            <i class="fas fa-exchange-alt"></i>
          </span>
        </div>

        <!-- destination language selection -->
        <div class="translator-dest-lang">
          <b-dropdown aria-role="list" id="dest-lang">
            <template #trigger="{}">
              <b-button :label="destLang" type="is-primary" />
            </template>

            <b-dropdown-item @click="changeDestLang('en')"
              >英文</b-dropdown-item
            >
            <b-dropdown-item @click="changeDestLang('zh')"
              >中文</b-dropdown-item
            >
            <b-dropdown-item @click="changeDestLang('jp')"
              >日文</b-dropdown-item
            >
          </b-dropdown>
        </div>

        <!-- "translate" button -->
        <div class="translator-btn">
          <b-button
            class="button"
            id="do-translate"
            @click="doTranslate"
            variant=""
            >翻&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;译</b-button
          >
        </div>
      </div>

      <!-- main translation container -->
      <div class="translate-box-container">
        <div class="translate-box">
          <textarea
            class=""
            type=""
            v-model="sourceText"
            :placeholder="sourceLangPlaceHolder"
            multiple
          />
          <div class="translate-toolbar">
            <button class="toolbar-btn"><i class="fas fa-copy"></i></button>
            <span class="toolbar-info">{{ srcWordCount }} / 2000 </span>
          </div>
        </div>

        <div class="translate-box disable-input">
          <textarea
            class=""
            type=""
            v-model="destText"
            multiple
            :placeholder="destLangPlaceHolder"
            disabled
          />
          <div class="translate-toolbar">
            <button class="toolbar-btn"><i class="fas fa-copy"></i></button>
          </div>
        </div>
      </div>

      <!-- footer -->
      <FooterBar />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';
// import { BIconArrowLeftRight } from 'bootstrap-vue';
import { listeners } from '../service.renderer.listener';
import { AsyncMainAPI } from '../service.renderer.manifest';

const engineMap = {
  baidu: '百度翻译',
  google: '谷歌翻译',
  youdao: '有道翻译',
};

const langMap = {
  zh: '中文',
  en: '英文',
  jp: '日语',
};
const placeHolderMap = {
  zh: '你好，很高兴认识你！',
  en: 'Hello, nice to meet you!',
  jp: 'こんにちは、はじめまして！',
};

export default Vue.extend({
  components: { Header, FooterBar },
  computed: {
    srcWordCount(): Number {
      return !this.sourceText ? 0 : this.sourceText.length;
    },
  },
  data() {
    return {
      selectedEngine: '谷歌翻译',
      engine: 'google',
      sourceLang: '中文',
      sourceLangCode: 'zh',
      destLang: '英文',
      destLangCode: 'en',
      sourceLangPlaceHolder: placeHolderMap['zh'],
      destLangPlaceHolder: placeHolderMap['en'],
      sourceText: '',
      destText: '',
    };
  },
  methods: {
    useEngine(engine: string) {
      if (engine !== 'baidu' && engine !== 'google' && engine != 'youdao') {
        alert('翻译引擎暂不支持');
        return;
      }
      console.log(engine, engineMap[engine]);
      this.selectedEngine = engineMap[engine];
      this.engine = engine;
    },
    changeSourceLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.sourceLang = langMap[lang];
      this.sourceLangCode = lang;
      this.sourceLangPlaceHolder = placeHolderMap[lang];
    },
    changeDestLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.destLang = langMap[lang];
      this.destLangCode = lang;
      this.destLangPlaceHolder = placeHolderMap[lang];
    },
    doTranslate() {
      if (!this.sourceText || this.sourceText == '') {
        return;
      }
      if (this.engine === 'baidu') {
        AsyncMainAPI.asyncBaiduTranslate({
          query: this.sourceText,
          from: this.sourceLangCode,
          to: this.destLangCode,
        });
      } else if (this.engine === 'google') {
        AsyncMainAPI.asyncGoogleTranslate({
          query: this.sourceText,
          from: this.sourceLangCode,
          to: this.destLangCode,
        });
      } else if (this.engine === 'youdao') {
        AsyncMainAPI.asyncYoudaoTranslate({
          query: this.sourceText,
          from: this.sourceLangCode,
          to: this.destLangCode,
        });
      }
    },
  },
  mounted() {
    listeners.onAsyncTranslate((event, arg) => {
      console.log('===== translate result ====')
      console.log(arg);
      if (arg && arg.engine === 'baidu') {
        if (arg.code === 0 && arg.data) {
          if (arg.data.trans_result && arg.data.trans_result.length > 0) {
            this.destText = arg.data.trans_result[0].dst;
          } else {
            this.destText = 'failed (' + arg.code + '): ' + arg.message;
          }
        }
      } else if (arg && arg.engine === 'google') {
        if (arg.data && arg.data.length > 0) {
          this.destText = arg.data;
        } else {
          this.destText = 'failed (' + arg.code + '): ' + arg.message;
        }
      } else if (arg && arg.engine === 'youdao') {
        if (arg.data && arg.data.length > 0) {
          this.destText = arg.data;
        } else {
          this.destText = 'failed (' + arg.code + '): ' + arg.message;
        }
      } else {
        // do nothing
        this.destText = 'engine not recognized';
      }
    });
  },
});
</script>

<style lang="scss" scoped>
.translate-container {
  display: flex;
  padding: 0;
  margin: 0;
}

.sidebar {
  padding: 0;
  margin: 0;
  height: 100%;
  width: 160px;
  // margin-right: 20px;
  display: block;
  // border-right: 1px solid #ccc;
  background: #ffffff;

  // 翻译引擎
  .translator-engine {
    width: 100%;
    display: block;
    background: #fff;
    border-bottom: 1px solid #f1f1f1;
    margin-bottom: 10px;

    &::v-deep .dropdown {
      .lang-icon {
        color: #4080eb;
      }
      .button {
        width: 100%;
        border: none;
        outline: none;
        display: block;
        margin: 0 auto;
        padding: 15px 0;
        background: #fff;
        font-size: 22px;
        color: #535353;
        cursor: pointer;

        &:hover {
          background: #fff;
        }
      }
      .background {
        display: none;
      }

      .dropdown-menu {
        width: 159px;
        background: #fff;
        position: fixed;
        .dropdown-content {
          display: flex;
          flex-direction: column;
          padding-top: 10px;
          padding-bottom: 10px;
          min-height: 200px;
          border-bottom: 1px solid #f1f1f1;

          .dropdown-item {
            width: 100%;
            text-align: center;
            font-size: 18px;
            padding-top: 2px;
            padding-bottom: 2px;
            color: #666;
            cursor: pointer;
            &:hover {
              color: #333;
            }
            &:active {
              color: #111;
            }
          }
        }
      }
    }
  }

  .translator-src-lang {
    &::v-deep .dropdown {
      .button {
        width: 100%;
        border: none;
        outline: none;
        display: block;
        margin: 0 auto;
        padding: 10px 0;
        background: #fff;
        font-size: 18px;
        font-weight: 700;
        cursor: pointer;

        &:hover {
          background: #fff;
        }
      }
      .background {
        display: none;
      }

      .dropdown-menu {
        width: 159px;
        background: #fff;
        position: fixed;
        .dropdown-content {
          display: flex;
          flex-direction: column;
          padding-top: 10px;
          padding-bottom: 10px;
          min-height: 200px;
          border-bottom: 1px solid #f1f1f1;
          .dropdown-item {
            width: 100%;
            text-align: center;
            font-size: 18px;
            padding-top: 2px;
            padding-bottom: 2px;
            color: #666;
            cursor: pointer;
            &:hover {
              color: #333;
            }
            &:active {
              color: #111;
            }
          }
        }
      }
    }
  }
  .translator-icon-container {
    width: 100%;
    display: flex;
    justify-content: center;
    padding: 10px 0;
    span {
      font-size: 18px;
      color: #4080eb;
    }
  }

  .translator-dest-lang {
    &::v-deep .dropdown {
      .button {
        width: 100%;
        border: none;
        outline: none;
        display: block;
        margin: 0 auto;
        padding: 10px 0;
        background: #fff;
        font-size: 18px;
        font-weight: 700;
        cursor: pointer;

        &:hover {
          background: #fff;
        }
      }
      .background {
        display: none;
      }

      .dropdown-menu {
        width: 159px;
        background: #fff;
        position: fixed;
        .dropdown-content {
          display: flex;
          flex-direction: column;
          padding-top: 10px;
          padding-bottom: 10px;
          min-height: 200px;
          border-bottom: 1px solid #f1f1f1;
          .dropdown-item {
            width: 100%;
            text-align: center;
            font-size: 18px;
            padding-top: 2px;
            padding-bottom: 2px;
            color: #666;
            cursor: pointer;
            &:hover {
              color: #333;
            }
            &:active {
              color: #111;
            }
          }
        }
      }
    }
  }

  // 翻译按钮
  .translator-btn {
    margin: 20px auto;
    .button {
      border: none;
      outline: none;
      display: block;
      margin: 0 auto;
      padding: 8px 26px;
      background: #4080eb;
      border-radius: 4px;
      box-shadow: 0 1px 6px 0 rgb(32 33 36 / 28%);
      font-size: 18px;
      color: #fff;
      cursor: pointer;
      &:hover {
        background: #2d65c7;
      }
      &:active {
        background: rgb(32, 83, 172);
      }
    }
  }
}

// 翻译内容区域
.translate-box-container {
  width: calc(100% - 160px);
  position: relative;
  background: #4080eb;

  .translate-box {
    height: 220px;
    width: calc(100% - 60px);
    margin: 20px 30px;
    border-radius: 12px;
    // border: 1px solid #ccc;
    padding: 6px;
    padding-top: 18px;
    background: #fff;
    box-shadow: 0 1px 6px 0 rgb(32 33 36 / 28%);

    textarea {
      margin: 0;
      padding-left: 10px;
      text-align: left;
      width: 100%;
      height: calc(100% - 38px);
      resize: none;
      border: none;
      background: #fff;
      font-size: 18px;
      color: #464646;

      &:focus {
        outline: none !important;
        border: none !important;
      }
      &:active {
        outline: none !important;
        border: none !important;
      }
      &:hover {
        outline: none !important;
        border: none !important;
      }
    }
    .translate-toolbar {
      height: 30px;
      width: 100%;
      display: flex;
      flex-direction: row-reverse;
      margin-bottom: 8px;
      .toolbar-btn {
        background: transparent;
        border: none;
        outline: none;
        width: 26px;
        height: 26px;
        margin: 2px 10px;
        font-size: 16px;
        line-height: 26px;
        color: #999;
        cursor: pointer;
        &:hover {
          color: #666;
        }
        &:active {
          color: #333;
        }
      }

      .toolbar-info {
        margin: 2px 10px;
        font-size: 14px;
        line-height: 22px;
        height: 26px;
        padding: 4px 2px;
        color: #999;
      }
    }

    .disable-input {
      background: #f9f9f9;
      textarea {
        background: #f9f9f9;
      }
    }
  }
}
</style>
