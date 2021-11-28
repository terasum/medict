import Vuex from 'vuex';
import { StoreDataType } from './StoreDataType';
import { listeners } from '../service.renderer.listener';
import { createByProc } from '@terasum/electron-call';

import DictAPI from '../..//worker/apis/DictAPI';

const stubByRenderer = createByProc('renderer', 'error');
const dictApi = stubByRenderer.use<DictAPI>('worker', 'DictAPI');



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
  currentWord: { dictid: '', word: '' }, // current user input word (in the search input box)
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
    updatePopHistory(state, word) { },
  },
  actions: {
    /**
     * 搜索单个输入词，返回一个列表
     * @param param0 state 状态提交
     * @param payload 在搜索框输入的搜索词 payload {dictid, word}
     */
    async asyncSearchWord({ commit, state }, payload) {
      // update current select word index to 0
      // 设置当前选择的词为 0
      commit('updateSelectedWordIdx', 0);
      // update current searching word
      // 更新当前的输入词为 payload
      commit('updateCurrentWord', payload);
      // invoke associate
      // 查询建议词条
      console.log('[RENDERER] asyncSearchWord', payload);
      const suggestResult = await dictApi.suggestWord(payload.dictid, payload.word);
      console.log('[RENDERER] asyncSearchWord result ', suggestResult);
      // 设置建议词列表
      commit('updateSuggestWords', suggestResult);
      if(suggestResult.length > 0) {
        setTimeout(() =>{
          this.dispatch('asyncFindWordPrecisly', 0)
        }, 600)
      }
    },
    // 打开自研文件夹
    asyncOpenDictResourceDir({commit, state}, {dictid}) {
      console.log('[RENDERER] store.index asyncOpenDictResourceDir', dictid)
    },
    asyncLoadResource({commit, state}, payload) {
      // TODO
      console.log('[RENDERER] store.index search resource', payload)
    },
    asyncUpdateSideBar(context, payload) {
      // TODO
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },
    /**
     * 点击边栏，将数据展示在主要框中
     * @param param0 状态参数
     * @param id 用户点击的词条id
     * @returns 设置当前显示的词条内容
     */
    async asyncFindWordPrecisly({ commit, state }, id) {
      if (state.suggestWords[id]) {
        commit('updateSelectedWordIdx', id);
        const selectWord = state.suggestWords[id];
        console.log('[RENDERER] findWordPrecisly', id, state.suggestWords[id])
        const wordDef = await dictApi.lookupWordPrecisely(selectWord.dictid, selectWord.keyText, selectWord.rofset)
        console.log('[RENDERER] findWordPrecisly def: ',wordDef)
        if (!wordDef){
          Store.commit('updateCurrentContent', Buffer.from('NOT Loaded').toString('base64'));
          return;
        }
        // 先展示未处理的
        const newContent = Buffer.from(wordDef!.definition).toString('base64');
        // 再展示已经后处理的
        const postWordDef = await dictApi.postHandleDef(selectWord.dictid, wordDef.keyText, wordDef.definition);
        console.log('[RENDERER] post handled', postWordDef);
        // 编码返回
        const postContent = Buffer.from((postWordDef).definition).toString('base64');
        Store.commit('updateCurrentContent', postContent);
      
      } else {
        console.error(`error on find word precisly ${state.suggestWords[id]}`);
      }
    },

    asyncAddNewDict({ state, commit }, dict) {
      // const result = dictAccessorApi.dictAddOne({ dict: dict });
      // TODO FIX
      // const dicts = dictAccessorApi.dictFindAll();
      // commit('updateDictionaries', dicts);
      // return result;
    },

    asyncDelNewDict({ state, commit }, dictid) {
      // const result = dictAccessorApi.dictDeleteOne({ dictid: dictid });
      // TODO FIX
      // const dicts = dictAccessorApi.dictFindAll();
      // commit('updateDictionaries', dicts);
      // return result;
    },
    postWebviewEvent({state, commit}, event) {
      console.log('[RENDERER] store.index webview catch event', event);
      if (!event || event.type !== 'ipc-message') {
        return;
      }
      if (event.channel === 'entryLinkWord') {
        console.log(`[RENDERER] store.index webview entryLinkWord event catched`);
        if(!event.args || event.args.length < 1) {
          return;
        }
        const entryLinkQuery = event.args[0];
        console.log('[RENDERER] store.index click Entry Link', entryLinkQuery);
        this.dispatch('asyncSearchWord', {dictid: entryLinkQuery.dictid, word: entryLinkQuery.keyText})
      }
      
    // // 通过event.channel的值来判断webview发送的事件名
    //   // @ts-ignore
    //   if (event.channel === 'onFindWordPrecisly') {
    //     console.log(`[async:mainWindow] response onFindWordPrecisly:`);
    //     console.log(event);
    //     const newContent = Buffer.from(
    //       // @ts-ignore
    //       event.args[0].definition,
    //       'utf8'
    //     ).toString('base64');
    //     this.$store.commit('updateCurrentContent', newContent);
    //     // @ts-ignore
    //     this.$store.commit('updateCurrentLookupWord', event.args[0].keyText);
    //   }

    //   // @ts-ignore


    }
  },
});

// setup default data asynchronizly
(function setupDefaultData() {
  let loadTicker = setInterval(async () => {
    // load all dictories first
    let dicts = await dictApi.loadAllIndexed();
    if (dicts) {
      console.log('[RENDERER] load dicts from worker', dicts);
      Store.commit("updateDictionaries", dicts);
      if (dicts.length > 0) {
        Store.commit('updateCurrentSelectDict', dicts[0]);
      }
      console.log('[RENDERER] cleanup load dictionary ticker');
      clearInterval(loadTicker);
    } else {
      console.log('[RENDERER] retry loading dictionaries...');
    }
  }, 1000);
})();



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
