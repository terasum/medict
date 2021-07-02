import Vuex from 'vuex';
import { MainProcAsyncAPI } from '../service.renderer.manifest';

const Store = new Vuex.Store({
  state: {
    // defaultWindow: '/preference',
    defaultWindow: '/',

    suggestWords: [],
    historyStack: [],
    currentWord: '',
    headerData: {
      // currentTab: '设置',
      currentTab: '词典',
    },
    sideBarData: {
      selectedWordIdx: 0,
      candidateWordNum: 0,
    },
    dictionaries: [
      {
        id: 0,
        dictIdGen: 'XiS1',
        dictIdCustom: 'eng-chinese',
        dictName: '英汉双解',
        dictMdxFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdx',
        dictMddFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdd',
        dictDescription: '英汉双解',
      },
      {
        id: 1,
        dictIdGen: 'XiS1',
        dictIdCustom: 'collins',
        dictName: '柯林斯大辞典',
        dictMdxFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdx',
        dictMddFilePath:
          '/Users/chenquan/Workspace/nodejs/medict/resources/eng-chinese.mdd',
        dictDescription: '柯林斯大辞典',
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
      // associate
      MainProcAsyncAPI.suggestWord(payload);
    },
    REFER_LINK_WORD({ commit, state }, word) {
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      commit('updateCurrentWord', word);
      MainProcAsyncAPI.findWordPrecisly(word);
    },
    ASYCN_UPDATE_SIDE_BAR(context, payload) {
      console.log(`async-dispatch ASYCN_UPDATE_SIDE_BAR ${payload}`);
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },
    FIND_WORD_PRECISLY({ commit, state }, id) {
      if (state.suggestWords[id]) {
        commit('updateSelectedWordIdx', id);
        MainProcAsyncAPI.findWordPrecisly(state.suggestWords[id]);
      } else {
        console.error(`error on find word precisly ${state.suggestWords[id]}`);
      }
    },
  },
});

export default Store;
