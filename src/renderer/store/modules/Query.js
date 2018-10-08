const state = {
  rawdef: '',
  samelist: {}
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
  }
}

export default {
  state,
  mutations,
  actions
}
