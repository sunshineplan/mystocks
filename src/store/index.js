import { createStore } from 'vuex'
import { post } from '@/misc.js'

export default createStore({
  state: {
    stock: {},
    line: [],
    chart: null,
    stocks: [],
    indices: {}
  },
  mutations: {
    stock(state, stock) { state.stock = stock },
    line(state, line) { state.line = line },
    chart(state, chart) { state.chart = chart },
    stocks(state, stocks) { state.stocks = stocks },
    indices(state, indices) { state.indices = indices }
  },
  actions: {
    async stock({ dispatch, commit }, payload) {
      const resp = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'realtime',
      });
      const stock = await resp.json();
      if (stock.name) {
        commit('stock', stock);
        dispatch('updateChart')
      }
    },
    async line({ dispatch, commit }, payload) {
      const resp = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'chart',
      });
      const json = await resp.json();
      if (json.chart) {
        commit('line', json.chart)
        dispatch('updateChart')
      }
    },
    async indices({ commit }) {
      const resp = await fetch('/indices')
      commit('indices', await resp.json())
    },
    updateChart({ commit, state }) {
      const chart = state.chart
      if (chart && state.stock.last) {
        chart.options.scales.yAxes[0].ticks.suggestedMin =
          state.stock.last / 1.01
        chart.options.scales.yAxes[0].ticks.suggestedMax =
          state.stock.last * 1.01
        chart.annotation.elements.PreviousClose.options.value = state.stock.last
      }
      if (chart && state.line.length && state.stock.now) {
        const data = state.line
        data[data.length - 1].y = state.stock.now
        chart.data.datasets[0].data = data
      }
      if (chart) chart.update()
      commit('chart', chart)
    },
    resetChart({ commit, state }) {
      commit('line', [])
      const chart = state.chart
      if (chart) chart.data.datasets[0].data = []
      if (chart) chart.update()
      commit('chart', chart)
    },
    destroyChart({ commit, state }) {
      if (state.chart) state.chart.update()
      commit('chart', null)
    }
  }
})
