const state = {
  state: 'normal',
  dict: {
    name: '21世纪大辞典'
  }
}

const mutations = {
  TRANS_TO_SEARCH (glb) {
    if (glb.state === 'normal') {
      glb.state = 'search'
    }
  },
  TRANS_TO_NORMAL (glb) {
    if (glb.state === 'search') {
      glb.state = 'normal'
    }
  }
}

const actions = {
  someAsyncTask ({ commit }) {
    // do something async
    commit('CHANGE_STATE')
  }
}

export default {
  state,
  mutations,
  actions
}
