<template>
  <div id="wrapper">
    <transition name="header-trans">
      <Header :global="global"/>
    </transition>
    <main>
    </main>
  </div>
</template>

<script>
  import SystemInformation from './LandingPage/SystemInformation'
  import Header from './Header.vue'

  export default {
    name: 'main-window',
    components: { SystemInformation, Header },
    data () {
      return {
        global: {
          /* search or normal */
          state: 'search'
          // state: 'normal'
        }
      }
    }
  }
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Monoton');
  #wrapper {
    background:#1377ED;
    height: 100vh;
    width: 100vw;
    padding-top: 22px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  
  main {
    display: flex;
    width: 410px;
  }

/*animations*/
.header-trans-enter-active {
  transition: all .3s ease;
}
.header-trans-leave-active {
  transition: all .8s cubic-bezier(1.0, 0.5, 0.8, 1.0);
}
.header-trans-enter, .header-trans-leave-to {
  transform: translateX(10px);
}
</style>
