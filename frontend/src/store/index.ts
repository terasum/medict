import Vuex from 'vuex';
import Vue from 'vue';
import { StoreDataType } from './StoreDataType';
import DictApi from '../apis/DictApi';
import StaticInfoAPI from '../apis/StaticInfoApi';
import { contents } from 'cheerio/lib/api/traversing';

// use vuex
Vue.use(Vuex);

const windowOpenApi = {
  openDictResourceDir: function(){console.error("error")}
};

const dictApi = {
  // lookupDefinition : function(a:any, b:any, c:any){console.error("error");},
  suggestWord: function(a:any,b:any){console.error("error");},
  postHandle: function(a:any, b:any, c:any) {console.error("error");},
  loadAllIndexed: function() {console.error("error");},
  getBaseDir: function(){console.error("error"); return "error"},
}

const newDictApi = new DictApi();
const staticInfoAPI = new StaticInfoAPI();

const state: StoreDataType = {
  // 首页 router 数据
  defaultWindow: '/',
  // 导航tab数据
  headerData: {
    currentTab: '词典',
  },

  // 边栏数据
  sideBarData: {
    selectedWordIdx: 0,
    candidateWordNum: 0,
  },

  // 当前词典列表
  dictionaries: [],

  // 当前词典搜索目录
  dictBaseDir: '',
  // 建议词列表
  suggestWords: [],
  // 历史栈
  historyStack: [],
  // 当前词 (包括词典id和词)
  currentWord: { dictid: '', word: '' }, // current user input word (in the search input box)
  // 当前搜索词
  currentLookupWord: '', // current actually search word
  // 当前实际选择词
  currentActualWord: '', // current actually return word
  // 当前展示词典内容
  currentContent: 'data:text/html;charset=utf-8,<html></html>',
  // 当前选择词典
  currentSelectDict: {
    id: "",
    name: "",
  },
  // 翻译API
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
    updateDictBaseDir(state, dirPath) {
      state.dictBaseDir = dirPath;
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
    async asyncSuggest({ commit, state }, payload) {
      // 设置当前选择的词为 0
      commit('updateSelectedWordIdx', 0);
      // 更新当前的输入词为 payload
      commit('updateCurrentWord', payload);
      // 查询建议词条
      console.log('[RENDERER] 查询建议词条 ', payload);
      const suggestResult = await newDictApi.suggest(payload.word);
      // const suggestResult = await dictApi.suggestWord(payload.dictid, payload.word);
      console.log('[RENDERER] 建议词条查询结果 ', suggestResult);
      if (suggestResult instanceof Error) {
        return
      }
      // 设置建议词列表
      for (let i = 0; i < suggestResult.length; i++){
        suggestResult[i].offset = i
      }
      commit('updateSuggestWords', suggestResult);
      // 暂时不默认选择第一个, 跳转的处理重新再写逻辑
    },


    /**
     * 点击边栏，将数据展示在主要框中
     * @param param0 状态参数
     * @param id 用户点击的词条id
     * @returns 设置当前显示的词条内容
     */
    async asyncLocateWord({ commit, state }, offset) {
      console.log("用户选定offset", offset);
      console.log("当前建议词列表", state.suggestWords);
      if (offset === state.sideBarData.selectedWordIdx) {
        return;
      }
      if (!state.suggestWords[offset]) {
        console.error(`[RENDERER] 精确定位词条失败,选定offset非法: ${offset}`);
        return;
      }

      commit('updateSelectedWordIdx', offset);
      const selectWord = state.suggestWords[offset];
      console.log(`[RENDERER] 精确定位词条: ${offset}`, state.suggestWords[offset]);
      // use another method
      let queryUrl = await staticInfoAPI.StaticSrvUrl()
      queryUrl += `?dict_id=${selectWord.dict_id}&raw_key_word=${selectWord.raw_key_word}&record_start=${selectWord.record_start}`
      console.log("queryurl: " +queryUrl);
      Store.commit('updateCurrentContent',queryUrl);

      // const wordDef = await newDictApi.lookupDefinition(selectWord.dict_id, selectWord.raw_key_word, selectWord.record_start)
      // if (wordDef instanceof Error) {
      //    console.error(`[RENDERER] 精确定位词条失败, 返回结果为错误`);
      //    console.error(wordDef);
      //   Store.commit('updateCurrentContent', Buffer.from('NOT Loaded').toString('base64'));
      //    return;

      // }
      // console.log("def 查询返回内容: ", wordDef);

      // // 先展示未处理的
      // const newContent = wordDef;

      // // 再展示已经后处理的
      // // TODO FIX
      // // const postWordDef = await dictApi.postHandle(selectWord.dictid, wordDef.keyText, wordDef.definition);
      // // console.debug('[RENDERER] 词条处理后结果', postWordDef);
      // // 编码返回
      // // TODO FIX
      // // const postContent = Buffer.from((postWordDef).definition).toString('base64');
      // Store.commit('updateCurrentContent', newContent);
    },

// 异步刷新词典数据
    async refreshDictData() {
      // do refresh immeiately 
      Store.commit('updateDictionaries', []);
      refresh();
      // 每15秒刷新一次
      let loadTicker = setInterval(refresh, 15000);

      async function refresh() {
        // load all dictories first
        let dicts = await newDictApi.dicts();
        if (!dicts) {
          console.log('[RENDERER] 当前词典为空，稍后尝试重新刷新');
          return;
        }

        console.log('[RENDERER] worker 加载词典成功', dicts);
        Store.commit("updateDictionaries", dicts);

        // 默认选择第一个词典
        if (dicts.length > 0) {
          Store.commit('updateCurrentSelectDict', dicts[0]);
        }

        // // 更新配置中的当前词典基本目录
        // let dictBaseDir = await dictApi.getBaseDir();
        // if (!dictBaseDir) {
        //   return;
        // } 

        // Store.commit("updateDictBaseDir", dictBaseDir)
        console.log('[RENDERER] 清除词典刷新定时器');
        clearInterval(loadTicker);
      }

    },

// 处理从 webview 返回的消息
    handleWebviewEvent({ state, commit }, event) {
      console.log('[RENDERER] store.index webview catch event', event);
      if (!event || event.type !== 'ipc-message') {
        return;
      }
      if (event.channel === 'entryLinkWord') {
        console.log(`[RENDERER] store.index webview entryLinkWord event catched`);
        if (!event.args || event.args.length < 1) {
          return;
        }
        const entryLinkQuery = event.args[0];
        console.log('[RENDERER] store.index click Entry Link', entryLinkQuery);
        this.dispatch('asyncSuggest', { dictid: entryLinkQuery.dictid, word: entryLinkQuery.keyText })
      }
    },

    // 打开资源文件夹
    asyncOpenDictResourceDir({ commit, state }, { dictid }) {
      // windowOpenApi.openDictResourceDir(dictid);
    },

    asyncLoadResource({ commit, state }, payload) {
      console.log('[RENDERER] store.index search resource', payload)
    },

    asyncUpdateSideBar(context, payload) {
      context.commit('updateCandidateWordNum', payload.candidateWordNum);
    },


    
  },
});

// 系统启动时设置默认数据
(function setupDefaultData() {
  // 立即刷新数据
  Store.dispatch('refreshDictData');
})();

export default Store;
