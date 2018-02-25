// import path from 'path'
import Mdict from 'js-mdict'
import path from 'path'
import cheerio from 'cheerio'
import fs from 'fs-extra-promise'

const state = {
  state: 'normal',
  dict: {
    name: 'Collins',
    path: path.join(__static, './dicts/Collins.mdx'),
    dictionary: null
  },
  searchWords: [],
  searchContent: '',
  styleContent: ''
}

const mutations = {
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
  initDict ({ commit, state }) {
    return commit('INIT_DICT')
  },
  searchWold ({ commit, state }, words) {
    if (words.length === 0) {
      commit('NOT_FOUND')
    }
    const word = words[words.length - 1]
    if (!word.k) {
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

export default {
  state,
  mutations,
  actions
}
