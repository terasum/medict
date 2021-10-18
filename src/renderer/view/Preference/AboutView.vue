<template>
  <div class="container settings-view">
    <div class="settings-title">
      <h3><i class="fas fa-info-circle"></i> APP 信息</h3>
    </div>
    <div class="settings-body">
      <section class="settings-section">
        <div class="section-title">
          <h4><i class="fas fa-rocket"></i> 软件信息</h4>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>软件简介</label>
              <p>Medict, 现代跨平台词典App</p>
            </div>
            <div class="input-container">
            <b-button
                class="btn"
                @click="onClickGithub('https://github.com/terasum/medict')"
                >访问 github</b-button
              >
            </div>
          </div>
          <div class="input-group">
            <div class="input-info">
              <label>软件版本</label>
              <p>v2.0.6</p>
            </div>
            <div class="input-container">
            </div>
          </div>
        </div>
      </section>

      <section class="settings-section">
        <div class="section-title">
          <h4><i class="fas fa-hands-helping"></i> 致谢</h4>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>开发人员</label>
              <p>
                <a href="javascript:void(0);" @click="onClickGithub('https://github.com/terasum')">Chen, Quan</a>
              </p>
            </div>
            <div class="input-container">
            </div>
          </div>

          <div class="input-group">
            <div class="input-info">
              <label>设计人员</label>
              <p>Zhang, Mingjiao</p>
            </div>
            <div class="input-container">
            </div>
          </div>
          <div class="input-group">
            <div class="input-info">
              <label>特别鸣谢</label>
              <p>Song, Xing</p>
            </div>
            <div class="input-container">
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
import { AsyncMainAPI } from '../../service.renderer.manifest';

export default Vue.extend({
  components: {},
  data(){
    return {
    contributors:[{
      id:1,
      imgSize: 45,
      nickname: 'Chen Quan',
      info:'Developer',
      avatarURL: 'https://avatars.githubusercontent.com/u/20042193?v=4',
    },
    {
      id:2,
      imgSize: 45,
      nickname: 'Song Xing',
      info:'PM',
      avatarURL: 'https://avatars.githubusercontent.com/u/38558937?v=4',
    },
    {
      id:3,
      imgSize: 45,
      nickname: 'Zhang Mingjiao',
      info:'Designer',
      avatarURL: 'https://tva1.sinaimg.cn/large/008i3skNly1gtkrf6clgoj60d00d0dgc02.jpg',
    },
      
    ]
    }
  },
  methods: {
    onClickGithub(url: string) {
      AsyncMainAPI.openUrlOnBrowser(url);
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



.contributors{
  margin-top: 20px;
}
.contributor {
  margin-top: 10px;
  display: flex;
  justify-content: center;
.contributor-user {
  margin: 5px 8px;
}
}
.contributor-infos{
  padding: 0;
  margin:0;
  border:none;
  background: #f1f1f1;
  user-select: none;

.nickname{
  display: block;
  text-align: center;
  color: #272727;
  font-weight: 400;
  margin:0;
  padding:0;
}

.info{
  display: block;
  text-align: center;
  color: #272727;
  font-weight: lighter;
  margin:0;
  padding:0;
}

}
</style>
