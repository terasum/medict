<style lang="scss" scoped>
@import '@/style/variables.scss';

.app-sidebar {
  display: flex;
  flex-direction: column;
  background-color: #fafafa;
  height: 100%;
  width: 100%;
  overflow: hidden;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  cursor:default;

  .sidebar-top {
    display: flex;
    padding: 0 10px;
    .logo {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: $layout-sidebar-logo-height;
      img {
        width: 100%;
      }
    }
    
  }
  .sidebar-content{
    padding: 4px 2px;
    margin: 10px 0 0 0;
    user-select: none;
    height: calc(100% -  $layout-sidebar-logo-height - $layout-footer-height);
    ul{
      height: 100%;
      overflow-y: scroll;
      list-style: none;
      margin: 0;
      padding: 0 0 0 6px;;
      li{
       margin: 0 4px 0 4px;
       padding: 0 0px 0 6px;;
       border-bottom: 1px solid #f1f1f1;
       border-radius: 3px;
       user-select: none;
       -webkit-user-select: none;
       &:hover{
          background-color: #f1f1f1;
          cursor: pointer;
       }
      }
    }

  }
}
</style>
<template>
  <div id="app-sidebar" class="app-sidebar">

    <div class="sidebar-top">
      <div class="logo">
        <img src="@/assets/images/logo.png"/>
      </div>
    </div>

    <div class="sidebar-content">
      <ul>
        <li v-for="item in dictQueryStore.queryPendingList" :key="item.id" @click="locate(item.id)">
          <span >{{ item.key_word }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import {useDictQueryStore} from '@/store/dict';
const dictQueryStore = useDictQueryStore();

function locate(entry_id) {
  dictQueryStore.locateWord(entry_id);
}


</script>
