<template>
    <div class="row header">
      <div class="header-navigate-btns">
        <button type="button" class="btn btn-light btn-nav btn-nav-left"><b-icon-chevron-compact-left /></button>
        <button type="button" class="btn btn-light btn-nav btn-nav-right"><b-icon-chevron-compact-right /></button>
      </div>
      <div class="header-search-box">
        <b-input-group>
          <template v-slot:prepend>
            <b-dropdown text="英汉双解" variant="info">
              <b-dropdown-item>柯林斯大辞典</b-dropdown-item>
              <b-dropdown-item>牛津大辞典</b-dropdown-item>
            </b-dropdown>
          </template>
          <b-form-input :disabled='displaySearchBox' v-model="searchWord" ></b-form-input>
          <b-button variant="info"><b-icon-search /></b-button>
        </b-input-group>
      </div>

      <div class="header-functions">
        <div class="fn-box" v-bind:class="{'fn-box-active':currentTab === '词典'}"  v-on:click='clickDictionary' >
          <span class="fn-box-icon">
            <b-icon-eye-fill />
          </span>
          <span class="fn-box-text">词典</span>
        </div>

        <div class="fn-box" v-bind:class="{'fn-box-active':currentTab === '插件'}"  v-on:click='clickPlugins' >
          <span class="fn-box-icon">
            <b-icon-grid />
          </span>
          <span class="fn-box-text">插件</span>
        </div>

        <div class="fn-box" v-bind:class="{'fn-box-active':currentTab === '设置'}" v-on:click='clickPreference' >
          <span class="fn-box-icon">
            <b-icon-gear-fill />
          </span>
          <span class="fn-box-text" >设置</span>
        </div>
      </div>
    </div>
</template>


<script lang="ts">
import Vue from "vue";
import Store from '../store/index';

interface HeaderComponentData extends Vue {
  searchWord: string,
  currentTab: string,
}

export default Vue.extend({
  props: {
    displaySearchBox: {
      type: Boolean,
      default: false,
    }
  },
  data() {
    return  {
      searchWord: "",
    }
  },
  computed: {
   count () {
      return (this.$store as typeof Store).state.count;
    },
    currentTab() {
       return (this.$store as typeof Store).state.headerData.currentTab;
    }
  },
  methods: {
     clickDictionary(event: any) {
      console.log(event);
      this.$store.commit('changeTab', '词典');

      if (this.$router.currentRoute.path !== '/') {
        this.$router.replace({ path: '/' });
      }
    },
    clickPlugins(event: any) {
      console.log(event);
      this.$store.commit('changeTab', '插件');

      if (this.$router.currentRoute.path !== '/plugins') {
        this.$router.replace({ path: '/plugins' });
      }
    },
    clickPreference(event: any) {
      console.log(event);
      this.$store.commit('changeTab', '设置');

      this.$store.commit('increment');
      console.log(this.$store.state.count);
      (this as HeaderComponentData).searchWord = this.$store.state.count  + "";
      console.log("更新router!");
      if (this.$router.currentRoute.path !== '/preference') {
        this.$router.replace({ path: '/preference' });
      }
    }
  }
});
</script>


<style lang="scss" scoped>
.header {
  height: 60px;
  background-color: #d84042;
  padding-top: 6px;
  .header-navigate-btns {
    height: 54px;
    max-width: 80px;
    padding: 0;
    margin: 0;
    .btn-nav {
      height: 26px;
      width: 26px;
      margin-top:12px;
      padding: 0;
      text-align: center;
      font-size: 12px;
      color: #fff;
      outline: none;
      border: #A63230 1px solid;
      background-color:#D84042;
      box-shadow: none;

      &:active {
        box-shadow: none;
        border: #333 1px solid;
        background-color: #D80034;
      }
      
      // release
       &:focus {
        outline: none;
        box-shadow: none;
        border: #A63230 1px solid;
        background-color:#D84042;
      }
    }
    .btn-nav-left {
      border-radius: 10px 0px 0px 10px;
      margin-left: 9px;
    }
   .btn-nav-right {
      border-radius: 0px 10px 10px 0px;
    }
  }
  .header-search-box {
    max-width: 360px;
    height: 54px;
    padding: 0;
    margin:0;
    padding-top:12px;
    &::v-deep{
      .btn-group, .btn-group-vertical {
        vertical-align: top;
      }
      // toggle button
      button:nth-child(1) {
        background-color: #fff;
        border: 1px solid #fff;
        border-radius: 20px 0px 0px 20px;
        height: 26px;
        font-size: 12px;
        line-height: 26px;
        padding:0;
        margin:0;
        padding-left: 10px;
        padding-right: 10px;
        box-shadow: none;
      }
      // search button
      button:nth-child(3){
        background-color:#fff;
        border: 1px solid #fff;
        border-radius: 0px 20px 20px 0px;
        height: 26px;
        padding:0;
        margin:0;
        padding-left: 10px;
        padding-right: 10px;
        box-shadow: none;
        line-height: 26px;
        font-size: 12px;
      }
      input{
        height: 26px;
        padding:0;
        margin:0;
        box-shadow: none;
        border: 1px solid #fff;
      }
      .form-control:disabled, .form-control[readonly] {
          background-color: #fff;
      }

    }
  }
  .header-functions {
    max-width: 306px;
    height: auto;
    padding: 0;
    margin: 0;
    margin-left:20px;
    display: flex;
    flex-direction: row;
    .fn-box-active {
        background-color: #BD3134;
    }
    .fn-box {
      width: 59px;
      height: 53px;
      margin-top: -6px;
      padding-top: 6px;
      cursor: pointer;
      border-radius: 5px;

      &:hover{
         background-color: #c73639;
      }
      .fn-box-icon {
        width: 26px;
        height: 26px;
        color: #f9dad9;
        display: block;
        font-size: 22px;
        line-height: 26px;
        margin-left: auto;
        margin-right: auto;
      }
      .fn-box-text {
        color: #f9dad9;
        margin-left: auto;
        margin-right: auto;
        text-align: center;
        width: 100%;
        display: block;
        font-size: 12px;
        user-select:none;
      }
    }
  }
}
</style>