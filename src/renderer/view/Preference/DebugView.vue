<template>
  <div class="container settings-view">
    <div class="settings-title">
      <h3><i class="fas fa-bug"></i> 开发者选项</h3>
    </div>
    <div class="settings-body">
      <section class="settings-section">
        <div class="section-title">
          <h4><i class="fas fa-file-alt"></i> 软件日志</h4>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>软件运行日志</label>
              <p>软件主进程日志，记录运行在后台的任务日志</p>
              <p>运行日志路径为: {{ loggerPath() }}</p>
            </div>
            <div class="input-container">
              <b-button
                class="btn"
                @click="onClickMainProcessLog"
                >打开运行日志</b-button
              >
            </div>
          </div>
        </div>
      </section>

      <section class="settings-section">
        <div class="section-title">
          <h4><i class="fas fa-terminal"></i> 开发者工具</h4>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>主APP开发工具</label>
              <p>可以通过以下按钮打开 chrome-dev-tool 面板(用于调试渲染进程)</p>
            </div>
            <div class="input-container">
              <b-button
                class="btn"
                @click="onClickDevBtn"
                >打开调试窗口</b-button
              >
            </div>
          </div>
        </div>
      </section>

      <section class="settings-section">
        <div class="section-title">
          <h4><i class="fas fa-archive"></i> 词典缓存</h4>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>词典资源缓存</label>
              <p>词典资源路径为: {{ resourcePath() }}</p>
            </div>
            <div class="input-container">
              <b-button
                class="btn"
                @click="onClickRescDir"
                >打开资源目录</b-button
              >
            </div>
          </div>
        </div>
      </section>

    </div>
    <div class="settings-footer">
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import {  SyncMainAPI } from '../../rpc.renderer.manifest';
import { createByProc } from '@terasum/electron-call';
import { WindowApi } from '../../../main/apis/WindowApi';
import {ipcRenderer} from 'electron';

const mainStub = createByProc('renderer');
const windowOpenApi = mainStub.use<WindowApi>('main', 'WindowApi');



import BugFill from '../../components/icons/bug-fill.icon.vue';
export default Vue.extend({
  components: { BugFill },
  methods: {
    saveConfig() {
      console.log(`save result`);
    },

    onClickDevBtn() {
      ipcRenderer.send('openDevTool');
    },

    onClickMainProcessLog() {
      windowOpenApi.openMainProcessLog();
    },
    onClickRescDir() {
      windowOpenApi.openResourceDir();
    },
    async loggerPath() {
      return SyncMainAPI.syncShowMainLoggerPath();
    },
    resourcePath() {
      return SyncMainAPI.syncGetResourceRootPath();
    },
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
  padding: 12px 0;

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
    min-width: 80px;
    &:active {
      background-color: #ddd;
      background-image: none;
    }
}
</style>
