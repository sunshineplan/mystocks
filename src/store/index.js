import { createStore } from 'vuex'
import Chart from 'chart.js'
import annotation from 'chartjs-plugin-annotation'
import { intraday } from '@/components/Chart.vue'
import { post } from '@/misc.js'

Chart.defaults.global.maintainAspectRatio = false
Chart.defaults.global.legend.display = false
Chart.defaults.global.hover.mode = 'index'
Chart.defaults.global.hover.intersect = false
Chart.defaults.global.tooltips.mode = 'index'
Chart.defaults.global.tooltips.intersect = false
Chart.defaults.global.tooltips.displayColors = false
Chart.defaults.global.animation.duration = 0
Chart.plugins.register({ annotation })

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
    async stock({ commit }, payload) {
      const resp = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'realtime',
      })
      const stock = await resp.json()
      if (stock.name) commit('stock', stock)
    },
    async line({ commit }, payload) {
      const resp = await post('/get', {
        index: payload.index,
        code: payload.code,
        q: 'chart',
      })
      const json = await resp.json()
      if (json.chart) commit('line', json.chart)
    },
    async indices({ commit }) {
      const resp = await fetch('/indices')
      commit('indices', await resp.json())
    },
    updateChart({ dispatch, commit, state }) {
      dispatch('destroyChart')
      const chart = new Chart(document.querySelector('#chart'), intraday)
      const stock = { ...state.stock }
      const line = [...state.line]
      if (stock.last) {
        chart.options.scales.yAxes[0].ticks.suggestedMin = stock.last / 1.01
        chart.options.scales.yAxes[0].ticks.suggestedMax = stock.last * 1.01
        chart.annotation.elements.PreviousClose.options.value = stock.last
      }
      if (line.length && stock.now) {
        line[line.length - 1].y = stock.now
        chart.config.data.datasets[0].data = line
      }
      chart.update()
      commit('chart', chart)
    },
    resetChart({ commit, state }) {
      commit('line', [])
      if (state.chart) {
        state.chart.destroy()
        commit('chart', new Chart(document.querySelector('#chart'), intraday))
      }
    },
    destroyChart({ commit, state }) {
      if (state.chart) state.chart.destroy()
      commit('chart', null)
    }
  }
})
