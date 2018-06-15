// import path from 'path'
import Mdict from 'js-mdict'
import path from 'path'
import cheerio from 'cheerio'
import fs from 'fs-extra-promise'
import ElectronStore from 'electron-store'

// const SET_CONTENT = '_SET_CONTENT'
const INIT_DICT = '_INIT_DICT'
const INIT_DICT_CSS = '_INIT_DICT_CSS'
const LOOK_UP = '_LOOK_UP'

const state = {
  store: new ElectronStore(),
  state: 'normal',
  searchWord: '',
  content: '',
  dict: {
    inited: false,
    name: 'Collins',
    path: path.join(__static, './dicts/Collins.mdx'),
    dictionary: null,
    dictionaryCss: {
      path: '',
      content: ''
    }
  },
  searchWords: [],
  searchContent: '',
  styleContent: ''
}

const mutations = {
  // 设置词典显示内容
  [LOOK_UP] (state, word) {
    state.searchWord = word
    let defination = state.dict.dictionary.lookup(word)
    console.log('mutations ' + defination)
    console.log('set content...')
    // if (state.dict.name === 'Collins') {
    const $ = cheerio.load(defination)
    $('link').remove()
    $('meta').remove()
    $('a[name=page_top]').remove()
    $('head').append('<style>' + state.dict.dictionaryCss.cssContent + '</style>')
    defination = $.html()
    // console.log(defination)
    // defination = defination.replace(/collinsEC.css/, () => {
    // return '~@/asstes/styles/CollinsEC.css'
    // })
    // }
    state.content = '<div id="medict-difinations" >' + defination + '</div><div id="medict-clear"></div>'
    // state.content = defination
  },
  [INIT_DICT] (state, dict) {
    console.log('mutations init dict')
    console.log(dict)
    state.dict.dictionary = dict
    state.dict.name = dict.attr().Title
    state.dict.inited = true
    console.log(state.dict)
  },
  [INIT_DICT_CSS] (state, cssContent, cssPath) {
    console.log('mutations init dict css')
    state.dict.dictionaryCss = {cssPath, cssContent}
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
      const cssPath = path.join(__static, './style/CollinsEC.css')
      const styleContent = fs.readFileSync(cssPath, 'utf8')
      commit(INIT_DICT_CSS, styleContent, cssPath)
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
