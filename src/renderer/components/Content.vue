<template>
  <div :class='["content", globalState.state]' v-html='content'>
  </div>  
</template>

<script>
export default {
  data () {
    return {
      globalState: this.$store.state.MedictGlobal
    }
  },
  computed: {
    content () {
      return this.globalState.searchContent ? this.globalState.searchContent : '<h3>NOT FOUND</h3>'
    }
  }
}
</script>


<style lang="scss" scoped>
.normal {
  display: none;
}

.search {
  width: 100vw;
  height: 100vh;
  display: flex;
  background: #fff;

  
}
#medict-difination{
    background-color: #fff;
    &::after{
      content: "";
      display: block;
      clear: both;
    }
  }
#medict-clear{
  clear: both;
}
</style>

