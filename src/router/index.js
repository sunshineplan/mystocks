import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import(/* webpackChunkName: "stock" */ '@/views/Stocks.vue') },
  { path: '/login', component: () => import(/* webpackChunkName: "login" */ '@/views/Login.vue') },
  { path: '/setting', component: () => import(/* webpackChunkName: "setting" */ '@/views/Setting.vue') },
  {
    name: 'stock',
    path: '/stock/:index/:code',
    component: () => import(/* webpackChunkName: "stock" */ '@/views/Stock.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
