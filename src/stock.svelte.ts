class MyStocks {
  username = $state('')
  component = $state('stocks')
  current = $state({ index: 'n/a', code: 'n/a' })
  refresh = $state(3)
  date = $state('')
  trading = $state(false)
}
export const mystocks = new MyStocks

export const info = async () => {
  const resp = await fetch('/info')
  const info = await resp.json()
  if (Object.keys(info).length) {
    mystocks.username = info.username
    mystocks.refresh = info.refresh
    if (info.date) {
      mystocks.date = info.date
      mystocks.trading = info.trading
    }
  }
}

class Toggler {
  status = $state(false)
  toggle() { this.status = !this.status }
}
export const isFlows = new Toggler
