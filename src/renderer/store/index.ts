import Vuex from 'vuex'
const Store = new Vuex.Store({
  state: {
    defaultWindow: '/preference',
    headerData:{
      currentTab: '设置',
    },
    count: 0
  },
  mutations: {
    increment (state) {
      state.count++
    },
    changeTab(state, tabName) {
      state.headerData.currentTab = tabName;
    }
  }
})

export default Store;