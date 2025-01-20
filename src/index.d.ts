interface Window {
  universal: string
  pubkey: string
}

interface Stock {
  index?: string
  code?: string
  name?: string
  now?: number
  change?: number
  percent?: string
  high?: number
  low?: number
  open?: number
  last?: number
  sell5?: { Price: number, Volume: number }[]
  buy5?: { Price: number, Volume: number }[]
}
