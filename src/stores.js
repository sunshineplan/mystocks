import { readable, writable } from 'svelte/store'
import Cookies from 'js-cookie'

export const stock = writable({})

export const refresh = readable(3, set => {
    const refresh = Cookies.get('Refresh')
    if (refresh) set(refresh)
})
