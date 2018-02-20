<template>
     <header :class="[ globalState.state === 'normal' ? 'center' : 'left' ]">
      <logo/>
      <search-box/>
      <history/>
      <select-dict/>
    </header> 
</template>



<script>
import SearchBox from './Header/SearchBox.vue'
import Logo from './Header/Logo.vue'
import History from './Header/History.vue'
import SelectDict from './Header/SelectDict.vue'

export default {
  name: 'Header',
  components: { SearchBox, Logo, History, SelectDict },
  data () {
    return {
      globalState: this.$store.state.MedictGlobal
    }
  }
}
</script>

<style lang="scss" scoped>
header {
  width: 100vw;
  overflow: hidden;
}
.center {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.left {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
}
</style>
