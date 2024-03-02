import { createRouter, createWebHashHistory } from 'vue-router';
import { IS_LOGGED_IN } from '@/stores/constants';
import type { RouteRecordName } from 'vue-router';

const routes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/Login.vue'),
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/pages/Home.vue'),
  },
  {
    path: '/signup',
    name: 'SignUp',
    component: () => import('@/pages/SignUp.vue'),
  },
];

const skippedProtectedEndpoints: RouteRecordName[] = ['Login', 'SignUp'];

const router = createRouter({
  // createWebHashHistory is for SPA to manage different states or views by using
  // hash in the URL for smooth navigation without reloading the entire page
  history: createWebHashHistory(),
  routes,
  // scrolling to the top of the screen
  scrollBehavior() {
    return {
      top: 0,
      left: 0,
      behavior: 'smooth',
    };
  },
});

// Navigation guard: https://router.vuejs.org/guide/advanced/navigation-guards.html
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem(IS_LOGGED_IN);

  // Added `to.name` to Avoid an infinite redirect
  if (isLoggedIn !== 'true' && !skippedProtectedEndpoints.includes(to.name!)) {
    next({ name: 'Login' });
  } else {
    next();
  }
});

export default router;
