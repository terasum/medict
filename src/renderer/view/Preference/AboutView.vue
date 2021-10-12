<template>
  <div class="container-fluid about-view">
    <div class="px-4 py-5 my-5 text-center">
      <h1 class="display-5 fw-bold">Medict</h1>
      <div class="col-lg-6 mx-auto">
        <p class="lead mb-4">现代跨平台词典App</p>
        <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
          <button
            type="button"
            class="btn btn-sm btn-outline-primary px-4 gap-3"
            @click="onClickGithub"
          >
            访问 github
          </button>
        </div>
      </div>

      <p class="contributors lead mb-4">Contributors</p>
      <div class="contributor">
        <div v-for="contributor in contributors" :key="contributor.id" class="contributor-user">
          <div class="avatar" :id="'contributor-' + contributor.id">
            <b-img :width="contributor.imgSize" :height="contributor.imgSize" rounded="circle"
              :src="contributor.avatarURL"></b-img>
          </div>
          <b-popover :target="'contributor-' + contributor.id" 
          triggers="hover" placement="bottom" boundary-paddin="0"
          custom-class="contributor-infos">
             <div class="nickname">{{ contributor.nickname }}<span></span></div>
             <div class="info"> {{ contributor.info }}<span></span></div>
          </b-popover>
        </div>


      </div>
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
    onClickGithub() {
      AsyncMainAPI.openUrlOnBrowser('https://github.com/terasum/medict');
    },
  },
});
</script>

<style lang="scss" scoped>
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
