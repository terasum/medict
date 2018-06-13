import Vue from 'vue'
import axios from 'axios'

import App from './App'
import router from './router'
import store from './store'
// import Antd from 'vue-antd-ui'
// import 'vue-antd-ui/dist/antd.css'
import iView from 'iview'
// import 'iview/dist/styles/iview.css'
import '@/assets/styles/iview-theme/medict.less'
import vuescroll from 'vuescroll'

// Vue.use(Antd)
Vue.use(iView)
Vue.use(vuescroll)

if (!process.env.IS_WEB) Vue.use(require('vue-electron'))

Vue.http = Vue.prototype.$http = axios
Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  components: { App },
  router,
  store,
  template: '<App/>',
  data: {
    ops: {
      vuescroll: {
        mode: 'native',
        // vuescroll's size(height/width) should be a percent(100%)
        // or be a number that is equal to its parentNode's width or
        // height ?
        sizeStrategy: 'percent',
        // pullRefresh or pushLoad is only for the slide mode...
        pullRefresh: {
          enable: false,
          tips: {
            deactive: 'Pull to Refresh',
            active: 'Release to Refresh',
            start: 'Refreshing...',
            beforeDeactive: 'Refresh Successfully!'
          }
        },
        pushLoad: {
          enable: false,
          tips: {
            deactive: 'Push to Load',
            active: 'Release to Load',
            start: 'Loading...',
            beforeDeactive: 'Load Successfully!'
          }
        },
        paging: false,
        zooming: true,
        snapping: {
          enable: false,
          width: 100,
          height: 100
        },
        // some scroller options
        scroller: {
          /** Enable bouncing (content can be slowly moved outside and jumps back after releasing) */
          bouncing: true,
          /** Enable locking to the main axis if user moves only slightly on one of them at start */
          locking: true,
          /** Minimum zoom level */
          minZoom: 0.5,
          /** Maximum zoom level */
          maxZoom: 3,
          /** Multiply or decrease scrolling speed **/
          speedMultiplier: 1,
          /** This configures the amount of change applied to deceleration when reaching boundaries  **/
          penetrationDeceleration: 0.03,
          /** This configures the amount of change applied to acceleration when reaching boundaries  **/
          penetrationAcceleration: 0.08,
          /** Whether call e.preventDefault event when sliding the content or not */
          preventDefault: true
        }
      },
      scrollPanel: {
        // when component mounted.. it will automatically scrolls.
        initialScrollY: false,
        initialScrollX: false,
        // feat: #11
        scrollingX: false,
        scrollingY: true,
        speed: 300,
        easing: undefined
      },
      //
      scrollContent: {
        // customize tag of scrollContent
        tag: 'div',
        padding: false,
        props: {},
        attrs: {}
      },
      //
      rail: {
        vRail: {
          width: '4px',
          pos: 'right',
          background: '#01a99a',
          opacity: 0
        },
        //
        hRail: {
          height: '4px',
          pos: 'bottom',
          background: '#01a99a',
          opacity: 0
        }
      },
      bar: {
        //
        vBar: {
          background: '#ccc',
          keepShow: false,
          opacity: 1,
          hover: false
        },
        //
        hBar: {
          background: '#aaa',
          keepShow: false,
          opacity: 1,
          hover: false
        }
      }
      // ...
    }
  }
}).$mount('#app')
