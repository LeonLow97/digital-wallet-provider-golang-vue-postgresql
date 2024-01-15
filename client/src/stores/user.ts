import { ref, watch, onBeforeMount } from 'vue';
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
  const storedUser = localStorage.getItem('user');
  if (storedUser) {
    user.value = JSON.parse(storedUser);
  }
  // Persist store throughout page reloads: https://github.com/vuejs/pinia/issues/309
  watch(
    user,
    (userVal) => {
      localStorage.setItem('user', JSON.stringify(userVal));
    },
    { deep: true }
  );

  const LOGIN_USER = (data: User) => {
    user.value.email = data.email;
    user.value.username = data.username;
    user.value.mobileNumber = data.mobileNumber;
    isLoggedIn.value = true;
  };

  const LOGOUT_USER = () => {
    isLoggedIn.value = false;
  };

  return {
    user,
    isLoggedIn,
    LOGIN_USER,
    LOGOUT_USER,
  };
});
