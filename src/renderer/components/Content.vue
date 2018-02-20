<template>
  <div :class='["content", globalState.state]'>

  </div>  
</template>

<script>
export default {
  data () {
    return {
      globalState: this.$store.state.MedictGlobal
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
</style>

