<template>
  <div class="container-fluid settings-view">
    <div class="section">
      <div class="section-title">百度翻译API设置</div>
      <div class="input-group">
        <label>APPID</label>
        <div class="input-container">
          <input
            type="text"
            class="form-control"
            placeholder="appid"
            v-model="baidu_appid"
          />
        </div>
      </div>
      <div class="input-group">
        <label>APPKEY</label>
        <div class="input-container">
          <input
            type="password"
            class="form-control"
            placeholder="appkey"
            v-model="baidu_appkey"
          />
        </div>
      </div>
    </div>
    <div class="btns-container">
      <div class="btns">
        <button
          class="btn btn-confirm btn-primary-outline"
          @click="saveConfig()"
        >
          确认
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { SyncMainAPI } from '../../service.renderer.manifest';

export default Vue.extend({
  components: {},
  data() {
    return {
      baidu_appkey: '',
      baidu_appid: '',
    };
  },
  methods: {
    saveConfig() {
      console.log(this.baidu_appkey);
      console.log(this.baidu_appid);
      const result = SyncMainAPI.saveTranslateBaiduApiConfig({
        appid: this.baidu_appid,
        appkey: this.baidu_appkey,
      });
      console.log(`save result ${result}`);
    },
  },
  mounted() {
    this.$nextTick(function () {
      const config = SyncMainAPI.loadTranslateApiConfig();
      if (config && config.hasOwnProperty('baidu')) {
        this.baidu_appkey = config.baidu.appkey;
        this.baidu_appid = config.baidu.appid;
      }
    });
  },
});
</script>

<style lang="scss" scoped>
.section {
  margin: 20px 10px;
  .input-group {
    margin-top: 10px;
    label {
      font-size: 14px;
      color: #666;
      min-width: 70px;
      margin-right: 10px;
      min-height: 30px;
    }
    .input-container {
      width: 15rem;
      min-height: 30px;
    }
  }
}
.btns-container {
  flex-direction: column-reverse;
  .btns {
    button {
      font-weight: 400;
      line-height: 1;
      border: 1px solid rgb(58, 140, 235);
    }
  }
}
</style>
