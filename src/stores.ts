import { writable } from 'svelte/store'

export const username = writable('')
export const component = writable('stocks')
export const current = writable({ index: 'n/a', code: 'n/a' })
export const refresh = writable(3)
export const date = writable('')
export const trading = writable(false)

export const info = async () => {
  const resp = await fetch('/info')
  const info = await resp.json()
  if (Object.keys(info).length) {
    username.set(info.username)
    refresh.set(info.refresh)
    if (info.date) {
      date.set(info.date)
      trading.set(info.trading)
    }
  }
}

const createFlows = () => {
  const { subscribe, update } = writable(false)
  return {
    subscribe,
    toggle: () => update(status => !status),
  }
}
export const flows = createFlows()
