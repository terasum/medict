import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'main',
      component: require('@/components/MainFrame').default,
      children: [
        {
          // 当 /setting 匹配成功，
          // Setting 会被渲染在 User 的 <router-view> 中
          path: 'setting',
          component: require('@/components/SettingFrame').default
        },
        {
          path: 'webview',
          name: 'webview',
          component: require('@/components/WebviewFrame').default
        }
      ]
    },
    {
      path: '*',
      redirect: '/webview'
    }
  ]
})
