<template>
 <div class="layout main-container">
    <div class="dragable" style="-webkit-app-region: drag"></div>
    <div class="main-content">
      <router-view></router-view>
    </div>        
 </div>

</template>

<style lang="scss" scoped>
  .main-container{
    background: #f6f6f6;
    height: 100%;
    width: 100%;
    overflow-y: hidden;
    
    .main-content{
      height: 100%;
      width: 100%;
    }
  }
  
.dragable{
  width: 100%;
  height: 22px;
  display: block;
  background: #1377ed ;
}
</style>


