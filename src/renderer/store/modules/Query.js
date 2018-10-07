const state = {
  rawdef: '',
  samelist: {}
}

const mutations = {
  UPDATE_RAW (state, raw) {
    state.rawdef = raw
  },
  UPDATE_LIST (state, list) {
    for (let i = 0; i < list.length; i++) {
      state.samelist[i] = list[i]
    }
  }
}

const actions = {
  updateRaw ({commit, state, getters}, raw) {
    console.log('store update the raw!!')
    commit('UPDATE_RAW', raw)
  },
  updateList ({commit, state, getters}, list) {
    console.log('store update the list!!')
    commit('UPDATE_LIST', list)
  }
}

export default {
  state,
  mutations,
  actions
}
