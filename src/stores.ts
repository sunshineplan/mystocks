import { writable } from 'svelte/store'

export const username = writable('')
export const component = writable('stocks')
export const current = writable({ index: 'n/a', code: 'n/a' })
export const refresh = writable(3)

const createFlows = () => {
  const { subscribe, update } = writable(false)
  return {
    subscribe,
    toggle: () => update(status => !status),
  }
}
export const flows = createFlows()
