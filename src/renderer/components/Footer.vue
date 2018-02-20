<template>
  <div :class='["footer", globalState.state]'>

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
  height: px;
  background: #1377ed;
}
</style>
