import { createStore } from 'vuex'
import { post } from '@/misc.js'

export default createStore({
  state: {
    stock: {},
    chart: [],
    stocks: [],
    indices: {}
  },
  mutations: {
    stock(state, stock) { state.stock = stock },
    chart(state, chart) { state.chart = chart },
    stocks(state, stocks) { state.stocks = stocks },
    indices(state, indices) { state.indices = indices }
  },
  actions: {
    async stock({ commit }, payload) {
      let response = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'realtime',
      });
      let stock = await response.json();
      if (stock.name) commit('stock', stock)
    },
    async chart({ commit }, payload) {
      let response = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'chart',
      });
      let json = await response.json();
      if (json.chart) commit('chart', json.chart)
    },
    async indices({ commit }) {
      const resp = await fetch('/indices')
      commit('indices', await resp.json())
    }
  }
})
