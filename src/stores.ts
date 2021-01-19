import { writable } from 'svelte/store'

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

export const username = writable('')
export const component = writable('stocks')
export const current = writable({ index: 'n/a', code: 'n/a' })
export const refresh = writable(3)
