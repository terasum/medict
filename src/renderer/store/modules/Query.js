import Store from 'electron-store'
const store = new Store()
const state = {
  rawdef: '',
  samelist: {},
  mdd: store.get('mdd'),
  mdx: store.get('mdx'),
  js: store.get('js'),
  css: store.get('css')
}

const mutations = {
  UPDATE_RAW (state, raw) {
    state.rawdef = raw
  },
  UPDATE_LIST (state, list) {
    console.log('store update the list!!')
    // 因为直接在原有对象修改不会触发更新，所以需要赋值新对象
    let newobj = {}
    for (let i = 0; i < list.length; i++) {
      newobj[i] = list[i]
    }
    // NOTE: must do this, if you modify samelist directory
    // the dom can't detect it
    state.samelist = Object.assign({}, newobj)
  },
  UPDATE_MDD (state, mddpath) {
    state.mdd = mddpath
    store.set('mdd', state.mdd)
  },
  UPDATE_MDX (state, mdxpath) {
    state.mdx = mdxpath
    store.set('mdx', state.mdx)
  },
  UPDATE_CSS (state, csspath) {
    state.css = csspath
    store.set('css', state.css)
  },
  UPDATE_JS (state, jspath) {
    state.js = jspath
    store.set('js', state.js)
  }
}

const actions = {
  updateRaw ({commit, state, getters}, raw) {
    console.log('store update the raw!!')
    commit('UPDATE_RAW', raw)
  },
  updateList ({commit, state, getters}, list) {
    console.log('store update the list!!')
    // commit('UPDATE_LIST', [])
    commit('UPDATE_LIST', list)
  },
  updateMDX ({commit}, mdxpath) {
    commit('UPDATE_MDX', mdxpath)
  },
  updateMDD ({commit}, mddpath) {
    commit('UPDATE_MDD', mddpath)
  },
  updateCSS ({commit}, csspath) {
    commit('UPDATE_CSS', csspath)
  },
  updateJS ({commit}, jspath) {
    commit('UPDATE_JS', jspath)
  }
}

export default {
  state,
  mutations,
  actions
}
