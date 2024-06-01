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
    component: () => import("@/pages/user/Login.vue"),
  },
  {
    path: "/home",
    name: "Home",
    component: () => import("@/pages/Home.vue"),
  },
  {
    path: "/signup",
    name: "SignUp",
    component: () => import("@/pages/user/SignUp.vue"),
  },
  {
    path: "/forgot-password",
    name: "ForgotPassword",
    component: () => import("@/pages/user/ForgotPassword.vue"),
  },
  {
    path: "/password-reset/:token",
    name: "PasswordReset",
    component: () => import("@/pages/user/PasswordReset.vue"),
  },
  {
    path: "/profile",
    name: "UserProfile",
    component: () => import("@/pages/user/UserProfile.vue"),
  },
  {
    path: "/settings",
    name: "Settings",
    component: () => import("@/pages/Settings.vue"),
  },
  {
    path: "/balances",
    name: "Balances",
    component: () => import("@/pages/balances/Balances.vue"),
  },
  {
    path: "/balances/:id",
    name: "Balance",
    component: () => import("@/pages/balances/Balance.vue"),
  },
  {
    path: "/transactions",
    name: "Transactions",
    component: () => import("@/pages/transactions/Transactions.vue"),
  },
  {
    path: "/beneficiary",
    name: "Beneficiary",
    component: () => import("@/pages/beneficiary/Beneficiary.vue"),
  },
  {
    path: "/wallets",
    name: "Wallets",
    component: () => import("@/pages/wallets/Wallets.vue"),
  },
  {
    path: "/wallet/:id",
    name: "Wallet",
    component: () => import("@/pages/wallets/Wallet.vue"),
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

const skippedProtectedEndpoints: RouteRecordName[] = [
  "Login",
  "SignUp",
  "ForgotPassword",
  "PasswordReset",
];

// Navigation guard: https://router.vuejs.org/guide/advanced/navigation-guards.html
router.beforeEach(async (to, from, next) => {
  const isLoggedIn = localStorage.getItem(IS_LOGGED_IN) === "true";

  // Skip calling GET_USER() if navigating to the login page
  if (skippedProtectedEndpoints.includes(to.name!)) {
    // If on the login page, check if the user is already authenticated
    if (to.name === "Login" && isLoggedIn) {
      try {
        const status = await GET_USER();
        if (status === 200) {
          next({ name: "Home" }); // Redirect to Home if authenticated
          return;
        }
      } catch (error) {
        localStorage.setItem(IS_LOGGED_IN, "false");
        console.error(error); // FOR DEVELOPMENT ONLY, REMOVE THIS!
      }
    }
    next();
    return;
  }

  // Redirect to login page if not logged in and not on a skipped protected endpoint
  if (!isLoggedIn && !skippedProtectedEndpoints.includes(to.name!)) {
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
        next({ name: "Error" });
        break;
    }
  } catch (error) {
    next({ name: "Login" });
  }
});

export default router;
