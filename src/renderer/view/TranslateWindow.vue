<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="translate-container" style="height: 100%">
      <div class="translate-btns">
        <div class="traslate-btn translator-engine">
          <b-dropdown aria-role="list"
            id="translator-engine"
          >
            <template #trigger="{ active }">
                <b-button
                    :label="selectedEngine"
                    type="is-primary"
                    :icon-right="active ? 'menu-up' : 'menu-down'" />
            </template>
            <b-dropdown-item @click="useEngine('baidu')">百度</b-dropdown-item>
            <b-dropdown-item @click="useEngine('google')">谷歌</b-dropdown-item>
            <b-dropdown-item @click="useEngine('bing')">必应</b-dropdown-item>
          </b-dropdown>

        </div>

        <div class="traslate-btn">
          <b-dropdown aria-role="list" id="src-lang">
             <template #trigger="{ active }">
                <b-button
                    :label="sourceLang"
                    type="is-primary"
                    :icon-right="active ? 'menu-up' : 'menu-down'" />
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
        <div class="traslate-btn">

          <b-dropdown aria-role="list" id="dest-lang">
             <template #trigger="{ active }">
                <b-button
                    :label="destLang"
                    type="is-primary"
                    :icon-right="active ? 'menu-up' : 'menu-down'" />
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
        <div class="traslate-btn">
          <b-button class="button" id="do-translate" @click="doTranslate" variant=""
            >翻&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;译</b-button
          >
        </div>
      </div>

      <div class="translate-box-container">
        <div class="translate-box-label"><span>源文本</span></div>
        <div class="translate-box">
          <textarea class="fullfill" type="" v-model="sourceText" multiple />
        </div>
        <div class="translate-box-label">
          <span>翻译</span>
        </div>
        <div class="translate-box">
          <textarea
            class="fullfill disable-input"
            type=""
            v-model="destText"
            multiple
            disabled
          />
        </div>
      </div>
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
  bing: '必应翻译',
};
const langMap = {
  zh: '中文',
  en: '英文',
  jp: '日语',
};

export default Vue.extend({
  components: { Header, FooterBar },
  data() {
    return {
      selectedEngine: '百度翻译',
      sourceLang: '中文',
      sourceLangCode: 'zh',
      destLang: '英文',
      destLangCode: 'en',
      sourceText: '',
      destText: '',
    };
  },
  methods: {
    useEngine(engine: string) {
      if (engine !== 'baidu') {
        alert('翻译引擎暂不支持');
        return;
      }
      console.log(engine, engineMap[engine]);
      this.selectedEngine = engineMap[engine];
    },
    changeSourceLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.sourceLang = langMap[lang];
      this.sourceLangCode = lang;
    },
    changeDestLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.destLang = langMap[lang];
      this.destLangCode = lang;
    },
    doTranslate() {
      if (!this.sourceText || this.sourceText == '') {
        return;
      }
      AsyncMainAPI.asyncBaiduTranslate({
        query: this.sourceText,
        from: this.sourceLangCode,
        to: this.destLangCode,
      });
    },
  },
  mounted() {
    listeners.onAsyncBaiduTranslate((event, arg) => {
      if (arg && arg.code === 0 && arg.data) {
        if (arg.data.trans_result && arg.data.trans_result.length > 0) {
          this.destText = arg.data.trans_result[0].dst;
        }
      } else {
        this.destText = '翻译失败: ' + arg.code;
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

.translate-btns {
  padding: 0;
  margin: 0;
  height: 100%;
  width: 160px;
  // margin-right: 20px;
  display: block;
  border-right: 1px solid #ccc;
  background: #f2f4f5;

  .traslate-btn {
    display: inline-block;
    height: 45px;
    width: 100%;
    line-height: 45px;
    padding-top: 0;
    .dropdown {
      padding: 0;
      margin: 0;
    }
  }

.traslate-btn::v-deep button {
    min-width: 120px;
    outline: none;
    border: none;
    background: #fff;
    line-height: 1rem;
    &:hover {
      background: #f1f1f1;
    }
    &:focus {
      outline: none;
      box-shadow: none;
    }
  }
  .dropdown::v-deep button {
    border: 1px solid #ccc;
    color: #333;
    &:hover {
      background: #f1f1f1;
    }
    &:focus {
      outline: none;
      box-shadow: none;
    }
  }

  .dropdown::v-deep ul {
    margin: 0;
    padding: 0;
  }
  .dropdown::v-deep ul > li {
    line-height: 2rem;
  }
  .dropdown::v-deep ul > li > a {
    margin: 0;
    padding: 0;
    padding-left: 0.5rem;
  }

  .translator-engine {
    width: 100%;
    border-bottom: 1px solid #ccc;
    display: block;
    #translator-engine{
     display: block; 
     margin-top: 10px;
     margin-bottom: 10px;
    }
    #translator-engine::v-deep button {
      width:100%;
      border:none;
      display: block;
      margin-left: auto;
      margin-right: auto;
      &:hover{
        background:#fff;
      }
    }
    &>ul{
      width:100%;
    }
  }

  #do-translate {
    background: #3a6bc7;
    color: #fff;
    margin-left: auto;
    margin-right: auto;
    display: block;
    margin-top: 20px;
    &:active {
      background: #2953a1;
      color: #f1f1f1;
    }
  }

  .dropdown::v-deep .dropdown-menu {
    min-width: 120px;
    max-width: 120px;
  }
  #src-lang{
    display: block;
  }
  #src-lang::v-deep button {
    width: 120px;
    text-align: center;
    margin-left:auto;
    margin-right:auto;
    margin-top: 15px;
    display: block;
  }
  #dest-lang{
    display: block;
  }
  #dest-lang::v-deep button {
    width: 120px;
    text-align: center;
    margin-left:auto;
    margin-right:auto;
    margin-top: 15px;
    display: block;
  }

  #dest-lang::v-deep .dropdown-menu {
    min-width: 120px;
    max-width: 120px;
  }

  #src-lang::v-deep .dropdown-menu {
    min-width: 120px;
    max-width: 120px;
    margin:0;
  }

  

  
}

.translate-box-container {
  width: calc(100% - 200px);
  position: relative;

  .translate-box-label {
    text-align: center;
    margin-top: 1rem;
    color: #6c6c6c;
  }

  .translate-box {
    height: 160px;
    width: 100%;
    padding: 0;
    margin: 0;

    textarea {
      margin: 0;
      padding-left: 10px;
      text-align: left;
      width: 100%;
      height: 100%;
      resize: none;
      border: none;
      border-radius: 3px;
      border: 1px solid #ccc;
      &:focus {
        outline: none !important;
        border: 1px solid #999;
      }
      &:active {
        border: 1px solid #999;
      }
      &:hover {
        border: 1px solid #aaa;
      }
    }

    .fullfill {
      width: 100%;
      height: 100%;
    }
    .disable-input {
      background: #f1f1f1;
    }
  }
}
</style>
