<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="translate-container" style="height: 100%">
      <div class="row translate-btns">
        <div class="traslate-btn translator-engine">
          <b-dropdown
            id="translator-engine"
            :text="selectedEngine"
            variant="outline-primary"
          >
            <b-dropdown-item @click="useEngine('baidu')">百度</b-dropdown-item>
            <b-dropdown-item @click="useEngine('google')">谷歌</b-dropdown-item>
            <b-dropdown-item @click="useEngine('bing')">必应</b-dropdown-item>
          </b-dropdown>
        </div>
        <div class="traslate-btn">
          <b-dropdown id="src-lang" :text="sourceLang" variant="outline-primary">
            <b-dropdown-item @click="changeSourceLang('EN')">英文</b-dropdown-item>
            <b-dropdown-item @click="changeSourceLang('ZH_CN')">中文</b-dropdown-item>
            <b-dropdown-item @click="changeSourceLang('JP')">日文</b-dropdown-item>
          </b-dropdown>
        </div>
        <div class="traslate-btn transfer-icon">
          <b-icon-arrow-left-right> </b-icon-arrow-left-right>
        </div>
        <div class="traslate-btn">
          <b-dropdown id="dest-lang" :text="destLang" variant="outline-primary">
            <b-dropdown-item @click="changeDestLang('EN')">英文</b-dropdown-item>
            <b-dropdown-item @click="changeDestLang('ZH_CN')">中文</b-dropdown-item>
            <b-dropdown-item @click="changeDestLang('JP')">日文</b-dropdown-item>
          </b-dropdown>
        </div>
        <div class="traslate-btn">
          <b-button id="do-translate" variant=""
            >翻&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;译</b-button
          >
        </div>
      </div>

      <div class="translate-box-label"><span>源文本</span></div>
      <div class="translate-box">
        <textarea class="fullfill" type="" value="在此输入文本" multiple />
      </div>
      <div class="translate-box-label"><span>翻译</span></div>
      <div class="translate-box">
        <textarea
          class="fullfill disable-input"
          type=""
          value="翻译"
          multiple
          disabled
        />
      </div>
      <FooterBar />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';
import { BIconArrowLeftRight } from 'bootstrap-vue';

const engineMap = {
  baidu: '百度翻译',
  google: '谷歌翻译',
  bing: '必应翻译',
};
const langMap = {
  ZH_CN: '中文',
  EN: '英文',
  JP: '日语',
};

export default Vue.extend({
  components: { Header, BIconArrowLeftRight, FooterBar },
  data() {
    return {
      selectedEngine: '百度翻译',
      sourceLang: '中文',
      destLang: '英文',
    };
  },
  methods: {
    useEngine(engine: string) {
      console.log(engine, engineMap[engine]);
      this.selectedEngine = engineMap[engine];
    },
    changeSourceLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.selectedEngine = langMap[lang];
    },
    changeDestLang(lang: string) {
      console.log(lang, langMap[lang]);
      this.selectedEngine = langMap[lang];
    },
  },
});
</script>

<style lang="scss" scoped>
.translate-container {
  display: block;
}

.translate-btns {
  padding: 0;
  margin: 0;
  margin-top: 15px;
  height: 45px;
  width: 100%;
  flex-wrap: nowrap;
  justify-content: space-between;
  align-items: center;
}

.traslate-btn {
  display: inline-block;
  height: 45px;
  width: auto;
  line-height: 45px;
  padding-top: 0;
  .dropdown {
    padding: 0;
    margin: 0;
  }
}
.transfer-icon {
  border: 1px solid red;
}
#do-translate {
  background: #d84042;
  color: #fff;
  &:active {
    background: #c73a3c;
    color: #f1f1f1;
  }
}

.dropdown::v-deep .dropdown-menu {
  min-width: 120px;
  max-width: 120px;
}

#src-lang::v-deep button {
  width: 120px;
  // text-align: left;
}

#dest-lang::v-deep button {
  width: 120px;
  // text-align: left;
}

#dest-lang::v-deep .dropdown-menu {
  min-width: 120px;
  max-width: 120px;
}

#src-lang::v-deep .dropdown-menu {
  min-width: 120px;
  max-width: 120px;
}

.traslate-btn::v-deep button {
  min-width: 120px;
  border: 1px solid #ccc;
  color: #333;
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
  min-width: 120px;
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
    &:focus {
      outline: none !important;
      border: 1px solid #ccc;
    }
    &:active {
      border: 1px solid #ccc;
    }
    &:hover {
      border: 1px solid #ccc;
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
</style>
