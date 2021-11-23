import Vuex from 'vuex';
import { StoreDataType } from './StoreDataType';
import { listeners } from '../service.renderer.listener';
import { createByProc } from '@terasum/electron-call';
import { DictAccessorApi } from '../../main/apis/DictAccessorApi';
import { WordQueryApi } from '../../main/apis/WordQueryApi';

const stubByRenderer = createByProc('renderer', 'error');
const dictAccessorApi = stubByRenderer.use<DictAccessorApi>('main', 'DictAccessorApi');
const wordQueryApi = stubByRenderer.use<WordQueryApi>('main', 'WordQueryApi');



const state: StoreDataType = {
  defaultWindow: '/',
  headerData: {
    currentTab: '词典',
  },

  sideBarData: {
    selectedWordIdx: 0,
    candidateWordNum: 0,
  },

  dictionaries: [],
  suggestWords: [],
  historyStack: [],
  currentWord: {dictid:'', word:''}, // current user input word (in the search input box)
  currentLookupWord: '', // current actually search word
  currentActualWord: '', // current actually return word
  currentContent: '',
  currentSelectDict: {
    id: "",
    alias: "",
    name: "",
  },
  translateApi: {
    baidu: {
      appid: '',
      appkey: '',
    },
  },
};

const Store = new Vuex.Store({
  state,
  mutations: {
    // update current word
    updateCurrentWord(state, word) {
      state.currentWord = word;
    },
    updateCurrentLookupWord(state, word) {
      state.currentLookupWord = word;
    },
    updateCurrentActualWord(state, word) {
      state.currentActualWord = word;
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
    updateDictionaries(state, dicts) {
      state.dictionaries = dicts;
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
      // update current select word index to 0
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      commit('updateCurrentWord', payload);
      // invoke associate
      // TODO FIX
      // wordQueryApi.suggestWord(payload);
    },
    asyncUpdateSideBar(context, payload) {
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },
    asyncFindWordPrecisly({ commit, state }, id) {
      if (state.suggestWords[id]) {
        commit('updateSelectedWordIdx', id);
        // TODO FIX
        // wordQueryApi.findWordPrecisly(state.suggestWords[id]);
      } else {
        console.error(`error on find word precisly ${state.suggestWords[id]}`);
      }
    },
    asyncAddNewDict({ state, commit }, dict) {
      const result = dictAccessorApi.dictAddOne({ dict: dict });
      // TODO FIX
      // const dicts = dictAccessorApi.dictFindAll();
      // commit('updateDictionaries', dicts);
      return result;
    },
    asyncDelNewDict({ state, commit }, dictid) {
      const result = dictAccessorApi.dictDeleteOne({ dictid: dictid });
      // TODO FIX
      // const dicts = dictAccessorApi.dictFindAll();
      // commit('updateDictionaries', dicts);
      return result;
    },
  },
});

function defaultSelectDict() {
  // TODO FIX
  // const dicts = dictAccessorApi.dictFindAll(undefined);
  // if (!dicts || dicts.length <= 0) {
  //   return { id: '', alias: '', name: '' };
  // }
  // return {
  //   id: dicts[0].id,
  //   alias: dicts[0].alias,
  //   name: dicts[0].name,
  // };
}

// setup default data asynchronizly
// TODO FIX
// (function setupDefaultData(){

//   let loadDicts = new Promise((resolve) =>{
//     // load all dictories first
//   let dicts = SyncMainAPI.dictFindAll(undefined);
//     resolve(dicts);
//   }).then(dicts =>{
//     Store.commit("updateDictionaries", dicts);
//   }).catch(err =>{
   
//   });

//   // current select dict
//   let loadCurrentSelect = new Promise((resolve) => {
//     resolve(defaultSelectDict());
//   }).then(dict =>{
//     Store.commit('updateCurrentSelectDict', dict);
//   });

// })();



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
    console.debug(`[store/index/listener]{onFindWordPrecisly}: resp:`);
    console.debug(args);
    const newContent = Buffer.from(args.definition, 'utf8').toString('base64');
    Store.commit('updateCurrentContent', newContent);
    Store.commit('updateCurrentLookupWord', args.sourceKey);
    Store.commit('updateCurrentActualWord', args.keyText);

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
    console.debug(`[store/index/listener]{onLoadDictResource}: resp:`);
    console.debug(args);
    // this.currentContent = args.definition.trim("\r\n\u0000");
  });
})();

export default Store;
