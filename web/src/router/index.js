// Composables
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/views/Transfer.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: () => import('@/views/Transactions.vue')
  },
  {
    path: '/transfer-funds',
    name: 'Transfer Funds',
    component: () => import('@/views/Transfer.vue')
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/HomeView.vue')
  }
]

export const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.name !== 'login' && localStorage.getItem('leon_access_token') === null) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router
