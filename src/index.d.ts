import type { ScatterDataPoint } from 'chart.js'

export interface Stock {
  index: string
  code: string
  name: string
  now: number
  change: number
  percent: string
  high: number
  low: number
  open: number
  last: number
  sell5: { Price: number, Volume: number }[]
  buy5: { Price: number, Volume: number }[]
  update?: string
}

export interface Flows {
  sector: string
  chart: ScatterDataPoint[]
}
