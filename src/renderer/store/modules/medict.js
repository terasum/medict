// import path from 'path'
import Mdict from 'js-mdict'
import path from 'path'
import cheerio from 'cheerio'
import fs from 'fs-extra-promise'
import ElectronStore from 'electron-store'

const SET_CONTENT = '_SET_CONTENT'
const INIT_DICT = '_INIT_DICT'
const INIT_DICT_CSS = '_INIT_DICT_CSS'
const LOOK_UP = '_LOOK_UP'

const state = {
  store: new ElectronStore(),
  state: 'normal',
  searchWord: 'a',
  content: '',
  dict: {
    inited: false,
    name: 'OPTED v.003',
    path: path.join(__static, './dicts/opted003.mdx'),
    dictionary: null,
    dictionaryCss: ''
  },
  searchWords: [],
  searchContent: '',
  styleContent: ''
}

const mutations = {
  // 设置词典显示内容
  [SET_CONTENT] (state, htmlContent) {
    console.log('mutations ' + htmlContent)
    state.content = htmlContent
  },
  [INIT_DICT] (state, dict) {
    console.log('mutations init dict')
    console.log(dict)
    state.dict.dictionary = dict
    state.dict.inited = true
  },
  [INIT_DICT_CSS] (state, cssContent) {
    console.log('mutations init dict css')
    state.dict.dictionaryCss = cssContent
  },
  [LOOK_UP] (state, word) {
    console.log(state.dict.dictionary)
    if (!state.dict.inited) {
      state.content = 'loadding...'
      return
    }
    state.content = state.dict.dictionary.lookup(word)
  },
  INIT_DICT2 (state) {
    // 判断词典是否已经初始化过，如果未初始化则进行初始化
    if (!state.store.get('dict_inited')) {
      console.log('this dict never running...')
    }
  },
  TRANS_TO_SEARCH (state, query) {
    if (state.state === 'normal') {
      state.state = 'search'
    }
  },
  TRANS_TO_NORMAL (state) {
    if (state.state === 'search') {
      state.state = 'normal'
    }
  },
  INIT_DICT (state) {
    state.dict.dictionary = new Mdict(state.dict.path).build()
    state.styleContent = fs.readFileSync(path.join(__static, './style/CollinsEC.css'), 'utf8')
    return state.dict.dictionary
  },
  SEARCH_WORD (state, defination) {
    if (state.dict.name === 'Collins') {
      const $ = cheerio.load(defination)
      $('link').remove()
      $('meta').remove()
      $('a[name=page_top]').remove()
      $('head').append('<style>' + state.styleContent + '</style>')
      defination = $.html()
      // defination = defination.replace(/collinsEC.css/, () => {
      // return '~@/asstes/styles/CollinsEC.css'
      // })
    }
    state.searchContent = '<div id="medict-difinations" >' + defination + '</div><div id="medict-clear"></div>'
  },
  NOT_FOUND (state) {
    state.searchContent = '<h3>NOT FOUND</h3>'
  }
}

const actions = {
  init ({dispatch, commit, state}) {
    return dispatch('loadDict').then((_dict) => {
      commit(INIT_DICT, _dict)
    })
  },
  loadDict ({commit, state, getters}) {
    if (!state.dict.inited) {
      console.log('loading...')
      const styleContent = fs.readFileSync(path.join(__static, './style/CollinsEC.css'), 'utf8')
      commit(INIT_DICT_CSS, styleContent)
      return new Mdict(state.dict.path).build()
    } else {
      console.log('dict loaded...')
    }
  },
  search ({commit, state, getters}, word) {
    // commit(SET_CONTENT, getters.userDataPath)
    if (!state.dict.inited) {
      console.log(' dict not inited')
      return
    }
    commit(LOOK_UP, word)
  },
  initDict ({ commit, state }) {
    return commit('INIT_DICT')
  },
  searchWold2 ({ commit, state }, words) {
    if (words.length === 0) {
      commit('NOT_FOUND')
    }
    const word = words[words.length - 1]
    if (!word || !word.k) {
      commit('NOT_FOUND')
    }
    if (state.dict.dictionary && word) {
      return state.dict.dictionary.then((_mdict) => {
        commit('SEARCH_WORD', _mdict.lookup(word.k))
      })
    }
  },
  predict ({ commit, state }) {
    // query the world
    if (state.dict.dictionary) {
      return state.dict.dictionary
    } else {
      return commit('INIT_DICT')
    }
  }

}

const getters = {
  userDataPath: state => {
    return state.store.cwd
  }
}

export default {
  state,
  mutations,
  actions,
  getters
}
