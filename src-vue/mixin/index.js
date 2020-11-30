import { color } from '@/misc.js'

export default {
  methods: {
    addColor(stock, val) {
      if (stock && stock.name != 'n/a') {
        switch (val) {
          case 'change':
          case 'percent':
            return color(stock.change)
          case 'now':
            return color(stock.last, stock.now)
          case 'high':
            return color(stock.last, stock.high)
          case 'low':
            return color(stock.last, stock.low)
          case 'open':
            return color(stock.last, stock.open)
        }
      }
    },
    gotoStock(stock) { this.$router.push(`/stock/${stock.index}/${stock.code}`) }
  }
}
