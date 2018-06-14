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
  template: '<App/>'
}).$mount('#app')
