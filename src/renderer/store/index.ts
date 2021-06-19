import Vuex from 'vuex'
const Store = new Vuex.Store({
  state: {
    headerData:{
      currentTab: '词典',
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