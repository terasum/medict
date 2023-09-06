<template>
  <div class="container settings-view">
    <div class="settings-title">
      <h3><i class="fas fa-language"></i> 翻译设置</h3>
    </div>
    <div class="settings-body">
      <section class="settings-section">
        <div class="section-title"><h4><i class="fas fa-binoculars"></i> 百度翻译</h4></div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>APPID</label>
              <p>在baidu 开放平台申请的 APPKID</p>
            </div>
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
            <div class="input-info">
              <label>APPKEY</label>
              <p>在baidu 开放平台申请的 APPKEY</p>
            </div>
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
      </section>

<section class="settings-section">
        <div class="section-title"><h4><i class="fas fa-lightbulb"></i> 有道翻译</h4></div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>APPID</label>
              <p> 有道智云申请的 APPKID</p>
            </div>
            <div class="input-container">
              <input
                type="text"
                class="form-control"
                placeholder="appid"
                v-model="youdao_appid"
              />
            </div>
          </div>
          <div class="input-group">
            <div class="input-info">
              <label>APPKEY</label>
              <p>在有道智云申请的 APPKEY</p>
            </div>
            <div class="input-container">
              <input
                type="password"
                class="form-control"
                placeholder="appkey"
                v-model="youdao_appkey"
              />
            </div>
          </div>
        </div>
      </section>

    </div>
    <div class="settings-footer">
      <button class="btn" @click="saveConfig()">确认</button>
    </div>
    <div class="row">
        <Alert 
        :message="alertMessage" 
        :status="alertStatus" 
        :alertActive="alertShow"
        v-on:close-alert="closeAlert"></Alert>
      </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Alert from '../../components/default/Alert.vue'

import { defineComponent } from 'vue';
export default defineComponent({
  components: {Alert},
  data() {
    return {
      baidu_appkey: '',
      baidu_appid: '',
      youdao_appkey: '',
      youdao_appid: '',

      alertMessage: '测试信息',
      alertStatus: 'info',
      alertShow: false,
    };
  },
  methods: {
    closeAlert() {
      this.alertShow = false;
      this.$emit('alert-closed');
    },
    showAlert(status: string, msg: string, callback?: ()=>any) {
      this.alertMessage = msg;
      this.alertStatus = status;
      this.alertShow = true;
      if(callback) {
        let _this = this;
        let _handler = function() {
          callback();
          _this.$off('alert-closed', _handler);
        }
        this.$on('alert-closed', _handler);
      }
    },
    saveConfig() {
      // TODO FIX
      // const baiduResult = SyncMainAPI.saveTranslateBaiduApiConfig({
      //   appid: this.baidu_appid,
      //   appkey: this.baidu_appkey,
      // });
      // console.log(`save result ${baiduResult}`);
      // const youdaoResult = SyncMainAPI.saveTranslateYoudaoApiConfig({
      //   appid: this.youdao_appid,
      //   appkey: this.youdao_appkey,
      // });
      // console.log(`save result ${youdaoResult}`);
      // this.showAlert('info', '保存成功!');
    },
  },
  mounted() {
    this.$nextTick(function () {
      // TODO FIX
      // const config = SyncMainAPI.loadTranslateApiConfig();
      // if (config && config.hasOwnProperty('baidu')) {
      //   this.baidu_appkey = config.baidu.appkey;
      //   this.baidu_appid = config.baidu.appid;
      // }
      // if (config && config.hasOwnProperty('youdao')) {
      //   this.youdao_appid = config.youdao.appid;
      //   this.youdao_appkey = config.youdao.appkey;
      // }
    });
  },
});
</script>

<style lang="scss" scoped>
.settings-title {
  padding: 10px 0 10px 12px;
  border-bottom: 1px solid #dfdfdf;
  h3 {
    font-size: 20px;
    font-weight: 700px;
    color: #777;
  }
}
.settings-section {
  display: flex;
  flex-direction: column;
  border-bottom: 1px solid #c1c1c3;
  margin: 0 10px;
  padding:12px 0;

  .section-title {
    width: 100%;
    padding: 2px 4px;
    font-size: 16px;
  }

  .section-body {
    padding-left: 16px;
  }

  .input-group {
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    .input-info {
      label {
        font-size: 15px;
        color: rgb(43, 103, 194);
      }
      p {
        font-size: 12px;
        color: #999;
      }
    }

    .input-container {
      display: flex;
      flex-direction: column;
      justify-content: center;
      max-width: 15rem;
      padding-right: 20px;
      input {
        height: 26px;
      }
    }
  }
}
.settings-body {
  min-height: 260px;
}
.settings-footer {
  display: flex;
  flex-direction: row-reverse;
  padding-right: 30px;
  margin-top: 10px;
  
}

.btn {
    background-color: rgb(252, 252, 252);
    background-image: linear-gradient(to bottom, #fcfcfc 0, #f1f1f1 100%);
    border-radius: 4px;
    border: 1px solid transparent;
    box-shadow: 0 1px 1px rgb(0 0 0 / 6%);
    border-color: #c2c0c2 #c2c0c2 #a19fa1;
    text-align: center;
    color: #666;
    font-size: 12px;
    display: inline-block;
    padding: 3px 8px;
    margin-bottom: 0;
    white-space: nowrap;
    height: 26px;
    width: 80px;
    &:active {
      background-color: #ddd;
      background-image: none;
    }
  }
</style>
