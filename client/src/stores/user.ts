import { ref, watch } from 'vue';
import { defineStore } from 'pinia';

interface User {
  email: string;
  username: string;
  mobileNumber: string;
}

export const useUserStore = defineStore('user', () => {
  // state
  const user = ref<User>({
    email: '',
    username: '',
    mobileNumber: '',
  });
  const isLoggedIn = ref(false);

  // initialize the store after the reactivity system is set up
  const storedUser = localStorage.getItem('USER_DETAILS');
  if (storedUser) {
    user.value = JSON.parse(storedUser);
  }

  const storedLoggedIn = localStorage.getItem('IS_LOGGED_IN');
  if (storedLoggedIn) {
    isLoggedIn.value = storedLoggedIn === 'true';
  }

  // Persist store throughout page reloads: https://github.com/vuejs/pinia/issues/309
  watch(
    [user, isLoggedIn],
    ([userVal, isLoggedInVal]) => {
      localStorage.setItem('USER_DETAILS', JSON.stringify(userVal));
      localStorage.setItem('IS_LOGGED_IN', isLoggedInVal.toString());
    },
    { deep: true }
  );

  const LOGIN_USER = (data: User) => {
    user.value = {
      email: data.email,
      username: data.username,
      mobileNumber: data.mobileNumber,
    };

    isLoggedIn.value = true;
  };

  const LOGOUT_USER = () => {
    user.value = {
      email: '',
      username: '',
      mobileNumber: '',
    };

    isLoggedIn.value = false;
  };

  return {
    user,
    isLoggedIn,
    LOGIN_USER,
    LOGOUT_USER,
  };
});
