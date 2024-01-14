import { createRouter, createWebHashHistory } from 'vue-router';

const routes = [
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

export default router;
