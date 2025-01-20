class Toggler {
  status = $state(false)
  toggle() { this.status = !this.status }
}

class stock {
  index = $state('')
  code = $state('')
  stock = $state<Stock>({})
  stared = $state(false)
  update = $state('')
  constructor(index?: string, code?: string) {
    if (index) this.index = index
    if (code) this.code = code
  }
  async goto() {
    window.history.pushState({}, '', `/stock/${this.index}/${this.code}`)
    const resp = await fetch('/star')
    if ((await resp.text()) == '1') this.stared = true
    else this.stared = false
  }
}

class MyStocks {
  username = $state('')
  component = $state('stocks')
  #toggler = new Toggler
  current = $state(new stock())
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
  goto(index: string, code: string): void
  goto(stock: Stock): void
  goto(a: string | Stock, b?: string) {
    if (typeof a === 'string') this.current = new stock(a, b)
    else this.current = new stock(a.index, a.code)
    this.component = 'stock'
    this.current.goto()
  }
}
export const mystocks = new MyStocks
