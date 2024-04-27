import { ref, watch } from "vue";
import { defineStore } from "pinia";
import type { User } from "../types/user";
import { USER_DETAILS, IS_LOGGED_IN } from "./constants";

export const useUserStore = defineStore("user", () => {
  // state
  const user = ref<User>({
    firstName: "",
    lastName: "",
    email: "",
    username: "",
    mobileNumber: "",
  });
  const isLoggedIn = ref(false);

  // initialize the store after the reactivity system is set up
  const storedUser = localStorage.getItem(USER_DETAILS);
  if (storedUser) {
    user.value = JSON.parse(storedUser);
  }

  const storedLoggedIn = localStorage.getItem(IS_LOGGED_IN);
  if (storedLoggedIn) {
    isLoggedIn.value = storedLoggedIn === "true";
  }

  // Persist store throughout page reloads: https://github.com/vuejs/pinia/issues/309
  watch(
    [user, isLoggedIn],
    ([userVal, isLoggedInVal]) => {
      localStorage.setItem(IS_LOGGED_IN, isLoggedInVal.toString());
      if (isLoggedInVal) {
        localStorage.setItem(USER_DETAILS, JSON.stringify(userVal));
      } else {
        localStorage.removeItem(USER_DETAILS);
      }
    },
    { deep: true },
  );

  const SAVE_USER = (data: User) => {
    user.value = {
      firstName: data.firstName,
      lastName: data.lastName,
      email: data.email,
      username: data.username,
      mobileNumber: data.mobileNumber,
    };
  };

  const LOGIN_USER = (data: User) => {
    user.value = {
      firstName: data.firstName,
      lastName: data.lastName,
      email: data.email,
      username: data.username,
      mobileNumber: data.mobileNumber,
    };

    isLoggedIn.value = true;
  };

  const LOGOUT_USER = () => {
    user.value = {
      firstName: "",
      lastName: "",
      email: "",
      username: "",
      mobileNumber: "",
    };

    isLoggedIn.value = false;
  };

  return {
    user,
    isLoggedIn,
    LOGIN_USER,
    LOGOUT_USER,
    SAVE_USER,
  };
});
