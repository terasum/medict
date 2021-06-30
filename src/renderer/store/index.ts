import Vuex from 'vuex';
import apis from '../../service/service.renderer.register';

const Store = new Vuex.Store({
  state: {
    defaultWindow: '/preference',
    suggestWords: [],
    historyStack: [],
    currentWord: '',
    headerData: {
      currentTab: '设置',
    },
    sideBarData: {
      selectedWordIdx: 0,
      candidateWordNum: 0,
    },
    dictionaries: [
      {
        id: 0,
        brief: '英汉双解',
        filename: 'eng-chinese.mdx',
        desc: '英汉双解词典',
        mddFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdx',
        mdxFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdd',
        cssFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.css',
        jsFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.js',
      },
      {
        id: 1,
        brief: '柯林斯',
        filename: 'collins.mdx',
        desc: '柯林斯大辞典',
        mddFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/collins.mdx',
        mdxFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/collins.mdd',
        cssFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/collins.css',
        jsFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/collins.js',
      },
    ],
    count: 0,
  },
  mutations: {
    updateCurrentWord(state, word) {
      state.currentWord = word;
    },
    increment(state) {
      state.count++;
    },
    changeTab(state, tabName) {
      state.headerData.currentTab = tabName;
    },
    updateCandidateWordNum(state, num) {
      state.sideBarData.candidateWordNum = num;
    },
    updateSelectedWordIdx(state, idx) {
      if (idx >= 0 && idx < state.sideBarData.candidateWordNum) {
        state.sideBarData.selectedWordIdx = idx;
      } else if (idx >= state.sideBarData.candidateWordNum) {
        state.sideBarData.selectedWordIdx =
          state.sideBarData.candidateWordNum - 1;
      } else if (idx < 0) {
        state.sideBarData.selectedWordIdx = 0;
      }
    },
    suggestWords(state, words) {
      state.suggestWords = words || [];
    },
    pushHistory(state, word) {
      if (state.historyStack.length > 512) {
      }
    },

    popHistory(state, word) {},
  },
  actions: {
    ASYNC_SEARCH_WORD({ commit, state }, payload) {
      console.log(`async-dispatch ASYNC_SEARCH_WORD ${payload}`);
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      commit('updateCurrentWord', payload);
      // do nothing
      apis['asyncSearchWord'](payload);
      // associate
      apis['suggestWord'](payload);
    },
    REFER_LINK_WORD({ commit, state }, word) {
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      commit('updateCurrentWord', word);
      apis['findWordPrecisly'](word);
    },
    ASYCN_UPDATE_SIDE_BAR(context, payload) {
      console.log(`async-dispatch ASYCN_UPDATE_SIDE_BAR ${payload}`);
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },
    FIND_WORD_PRECISLY({ commit, state }, id) {
      if (state.suggestWords[id]) {
        commit('updateSelectedWordIdx', id);
        apis['findWordPrecisly'](state.suggestWords[id]);
      } else {
        console.error(`error on find word precisly ${state.suggestWords[id]}`);
      }
    },
  },
});

export default Store;
