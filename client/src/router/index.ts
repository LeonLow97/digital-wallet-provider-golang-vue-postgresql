import { createRouter, createWebHashHistory } from "vue-router";
import { IS_LOGGED_IN } from "@/stores/constants";
import type { RouteRecordName } from "vue-router";
import { GET_USER } from "@/api/user";

const routes = [
  {
    path: "/",
    redirect: "/home",
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/pages/Login.vue"),
  },
  {
    path: "/home",
    name: "Home",
    component: () => import("@/pages/Home.vue"),
  },
  {
    path: "/signup",
    name: "SignUp",
    component: () => import("@/pages/SignUp.vue"),
  },
  {
    path: "/profile",
    name: "UserProfile",
    component: () => import("@/pages/UserProfile.vue"),
  },
  {
    path: "/settings",
    name: "Settings",
    component: () => import("@/pages/Settings.vue"),
  },
  {
    path: "/balances",
    name: "Balances",
    component: () => import("@/pages/Balances.vue"),
  },
  {
    path: "/transactions",
    name: "Transactions",
    component: () => import("@/pages/Transactions.vue"),
  },
  {
    path: "/transfer",
    name: "Transfer",
    component: () => import("@/pages/Transfer.vue"),
  },
  {
    path: "/wallets",
    name: "Wallets",
    component: () => import("@/pages/Wallets.vue"),
  },
  {
    path: "/error",
    name: "Error",
    component: () => import("@/pages/Error.vue"),
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
      behavior: "smooth",
    };
  },
});

const skippedProtectedEndpoints: RouteRecordName[] = ["Login", "SignUp"];
// Navigation guard: https://router.vuejs.org/guide/advanced/navigation-guards.html
router.beforeEach(async (to, from, next) => {
  const isLoggedIn = localStorage.getItem(IS_LOGGED_IN);

  // Skip calling GET_USER() if navigating to the login page
  if (skippedProtectedEndpoints.includes(to.name!)) {
    next();
    return;
  }

  // Redirect to login page if not logged in and not on a skipped protected endpoint
  if (isLoggedIn !== "true" && !skippedProtectedEndpoints.includes(to.name!)) {
    next({ name: "Login" });
    return;
  }

  try {
    // Call GET_USER() to refresh session
    const status = await GET_USER();
    switch (status) {
      case 200:
        next();
        break;
      case 401:
        localStorage.setItem(IS_LOGGED_IN, "false");
        next({ name: "Login" });
        break;
      default:
        localStorage.setItem(IS_LOGGED_IN, "false");
        next({ name: "Error" }); // Redirect to an error page
        break;
    }
  } catch (error) {
    // Handle error refreshing session
    console.error("Error refreshing session:", error);
    next({ name: "Login" });
  }
});

export default router;
