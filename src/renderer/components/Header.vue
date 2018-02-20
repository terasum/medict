<template>
     <header :class="[ global.state === 'normal' ? 'center' : 'left' ]">
      <logo :global="global" />
      <history :global="global" />
      <search-box :global="global" />
      <select-dict :global="global" />
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
    props: {
      global: {
        type: Object,
        required: true
      }
    }
}
</script>

<style lang="scss" scoped>
  header{
    width: 100vw;
    overflow: hidden;
  }
  .center {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .left{
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
  }
</style>
