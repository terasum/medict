import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'main',
      component: require('@/components/MainFrame').default
    },
    {
      'path': 'content/',
      name: 'content',
      component: require('@/components/ContentFrame').default
    },
    {
      'path': 'background/',
      name: 'background',
      component: require('@/components/Background').default
    },

    {
      path: '*',
      redirect: '/'
    }
  ]
})
