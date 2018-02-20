<template>
  <div id="wrapper">
    <transition name="header-trans">
      <Header :global="global"  />
    </transition>
    <main :class='[global.state === "normal" ? "hide" : "" ]'>
      <Content :global="global"/>
    </main>
    <Footer :global="global"></Footer>
  </div>
</template>

<script>
import SystemInformation from './LandingPage/SystemInformation'
import Header from './Header.vue'
import Content from './Content.vue'
import Footer from './Footer.vue'

export default {
  name: 'main-window',
  components: { SystemInformation, Header, Content, Footer },
  data () {
    return {
      global: this.$store.state.MedictGlobal
    }
  },
  methods: {
    transState () {
      this.$store.commit('CHANGE_STATE')
    }
  }
}
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Monoton");
#wrapper {
  background: #1377ed;
  height: 100vh;
  width: 100vw;
  padding-top: 22px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

main {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  margin-top: 5px;
}

.hide {
  display: none;
}

/*animations*/
.header-trans-enter-active {
  transition: all 0.4s ease;
}
.header-trans-leave-active {
  transition: all 0.8s cubic-bezier(1, 0.5, 0.8, 1);
}
.header-trans-enter,
.header-trans-leave-to {
  transition: transform 2s;
}
</style>
