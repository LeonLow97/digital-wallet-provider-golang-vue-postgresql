import { ref } from 'vue';
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

  const LOGIN_USER = (data: User) => {
    user.value.email = data.email;
    user.value.username = data.username;
    user.value.mobileNumber = data.mobileNumber;
    isLoggedIn.value = true;
  };

  return {
    user,
    isLoggedIn,
    LOGIN_USER,
  };
});
