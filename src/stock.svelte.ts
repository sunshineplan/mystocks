class Toggler {
  status = $state(false)
  toggle() { this.status = !this.status }
}

class MyStocks {
  username = $state('')
  component = $state('stocks')
  #toggler = new Toggler
  current = $state({ index: 'n/a', code: 'n/a' })
  refresh = $state(3)
  date = $state('')
  trading = $state(false)
  async info() {
    const resp = await fetch('/info')
    const info = await resp.json()
    if (Object.keys(info).length) {
      this.username = info.username
      this.refresh = info.refresh
      if (info.date) {
        this.date = info.date
        this.trading = info.trading
      }
    }
  }
  get isFlows() {
    return this.#toggler.status
  }
  toggle() {
    this.#toggler.toggle()
  }
}
export const mystocks = new MyStocks
