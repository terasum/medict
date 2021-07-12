import Vuex from 'vuex';
import { AsyncMainAPI, SyncMainAPI } from '../service.renderer.manifest';
import { StoreDataType } from './StoreDataType';
import { listeners } from '../service.renderer.listener';

function defaultSelectDict() {
  const dicts = SyncMainAPI.dictFindAll(undefined);
  if (!dicts || dicts.length <= 0) {
    return { id: '', alias: '', name: '' };
  }
  return {
    id: dicts[0].id,
    alias: dicts[0].alias,
    name: dicts[0].name,
  };
}

const state: StoreDataType = {
  defaultWindow: '/preference',
  // defaultWindow: '/',
  headerData: {
    currentTab: '设置',
    // currentTab: '词典',
  },

  sideBarData: {
    selectedWordIdx: 0,
    candidateWordNum: 0,
  },

  dictionaries: SyncMainAPI.dictFindAll(undefined),
  suggestWords: [],
  historyStack: [],
  currentWord: '',
  currentContent: '',
  currentSelectDict: defaultSelectDict(),
};

const Store = new Vuex.Store({
  state,
  mutations: {
    // update current word
    updateCurrentWord(state, word) {
      state.currentWord = word;
    },
    updateCurrentContent(state, content) {
      state.currentContent = content;
    },
    updateCurrentSelectDict(state, dict) {
      if (!dict || !dict.id || dict.id === '') {
        return;
      }
      state.currentSelectDict.id = dict.id;
      state.currentSelectDict.alias = dict.alias || '';
      state.currentSelectDict.name = dict.name || '';
    },
    updateTab(state, tabName) {
      state.headerData.currentTab = tabName;
    },
    updateCandidateWordNum(state, num) {
      state.sideBarData.candidateWordNum = num;
    },
    updateDictionaries(state) {
      state.dictionaries = SyncMainAPI.dictFindAll(undefined);
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
    updateSuggestWords(state, words) {
      state.suggestWords = words || [];
    },
    updatePushHistory(state, word) {
      if (state.historyStack.length > 512) {
      }
    },
    updatePopHistory(state, word) {},
  },
  actions: {
    asyncSearchWord({ commit, state }, payload) {
      console.log(`async-dispatch asyncSearchWord ${payload}`);
      // update current select word index to 0
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      commit('updateCurrentWord', payload);
      // invoke associate
      AsyncMainAPI.suggestWord(payload);
    },
    asyncUpdateSideBar(context, payload) {
      console.log(`async-dispatch asyncUpdateSideBar ${payload}`);
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },
    asyncFindWordPrecisly({ commit, state }, id) {
      if (state.suggestWords[id]) {
        commit('updateSelectedWordIdx', id);
        AsyncMainAPI.findWordPrecisly(state.suggestWords[id]);
      } else {
        console.error(`error on find word precisly ${state.suggestWords[id]}`);
      }
    },
    asyncAddNewDict({ state, commit }, dict) {
      const result = SyncMainAPI.dictAddOne({ dict: dict });
      commit('updateDictionaries');
      return result;
    },
    asyncDelNewDict({ state, commit }, dictid) {
      const result = SyncMainAPI.dictDeleteOne({ dictid: dictid });
      commit('updateDictionaries');
      return result;
    },
  },
});

(function setupListener() {
  /**
   * onSuggestWord catch main-process return suggest word response
   * @param event: event source, main-process
   * @param args: suggestion word payload,
   * [{
   *   dictid: "ar3e0x"
   *   id: 0
   *   keyText: "wack"
   *   rofset: 156414748
   * }]
   */
  listeners.onSuggestWord((event: any, args: any) => {
    console.log(`[store/index/listener]{onSuggestWord}: resp:`);
    console.log(args);
    Store.dispatch('asyncUpdateSideBar', {
      candidateWordNum: args.length,
    });
    Store.commit('updateSuggestWords', args);
  });

  /**
   * onFindWordPrecisly catch onFindWordPrecisly event, render main definition of word
   * @param event: event source, main-process
   * @param args: word definition payload,
   * {
   *   definition: "<html>"
   * }
   */
  listeners.onFindWordPrecisly((event: any, args: any) => {
    console.log(`[store/index/listener]{onFindWordPrecisly}: resp:`);
    console.log(args);
    const newContent = Buffer.from(args.definition, 'utf8').toString('base64');
    Store.commit('updateCurrentContent', newContent);
  });

  /**
   * onLoadDictResource catch onLoadDictResource event, render resource definition
   * @param event: event source, main-process
   * @param args: word definition payload,
   *                              * {
   *   definition: "<bbn>"
   * }
   */
  listeners.onLoadDictResource((event: any, args: any) => {
    console.log(`[store/index/listener]{onLoadDictResource}: resp:`);
    console.log(args);
    // this.currentContent = args.definition.trim("\r\n\u0000");
  });
})();

export default Store;
