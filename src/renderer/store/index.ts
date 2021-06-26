import Vuex from 'vuex';
import apis from "../../service/service.renderer.register";

const Store = new Vuex.Store({
  state: {
    defaultWindow: '/',
    headerData:{
      currentTab: '词典',
    },
    dictionaries: [
      {
        id: 0,
        brief: "英汉双解",
        filename: "eng-chinese.mdx",
        desc: "英汉双解词典",
        mddFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdx",
        mdxFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdd",
        cssFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.css",
        jsFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.js",
      },
      {
        id: 1,
        brief: "柯林斯",
        filename: "collins.mdx",
        desc: "柯林斯大辞典",
        mddFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/collins.mdx",
        mdxFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/collins.mdd",
        cssFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/collins.css",
        jsFilePath:
          "/Users/chenquan/Workspace/nodejs/medict/resources/collins.js",
      },
    ],
    count: 0
  },
  mutations: {
    increment (state) {
      state.count++
    },
    changeTab(state, tabName) {
      state.headerData.currentTab = tabName;
    }
  },
  actions: {
    ASYNC_SEARCH_WORD(state, payload) {
      console.log(`async-dispatch ${payload.word}`);
      apis['asyncSearchWord'](payload);
      apis['suggestWord'](payload);
    }
  }
})

export default Store;