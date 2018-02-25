<template>
  <div id="wrapper" :class="[initd ? '': 'hide']" @did-finish-load="init">
    <div class="dragable" style="-webkit-app-region: drag"></div>
    <Header :global="global"  />
    <main :class='[global.state === "normal" ? "hide" : "" ]'>
      <Content :global="global"/>
    </main>
    <Footer :global="global"></Footer>
  </div>
</template>

<script>
import SystemInformation from './Attachments/SystemInformation'
import Header from './Header.vue'
import Content from './Content.vue'
import Footer from './Footer.vue'

export default {
  name: 'main-window',
  components: { SystemInformation, Header, Content, Footer },
  data () {
    return {
      global: this.$store.state.MedictGlobal,
      inited: false
    }
  },
  computed: {
    initd () {
      return this.inited
    }
  },
  methods: {
    transState () {
      this.$store.commit('CHANGE_STATE')
    },
    init () {
      this.$store.dispatch('initDict').then((_mdict) => {
        this.inited = true
      })
    }
  },
  mounted: function () {
    this.$nextTick(function () {
      this.init()
      // Code that will run only after the
      // entire view has been rendered
    })
  }
}
</script>

<style lang="scss" scoped>
@import url("https://fonts.googleapis.com/css?family=Monoton");
.dragable{
  width: 100%;
  height: 22px;
  display: block;
  background: #1377ed ;
}

#wrapper {
  background: #1377ed;
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  align-items: center;
}

main {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow-x: hidden;
  background: #fff;
  margin-top: 5px;
  padding: 10px 40px;
}

.hide {
  display: none;
}

</style>
