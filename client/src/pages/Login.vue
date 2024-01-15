<template>
  <div>
    <h1>Login Page</h1>
    <form @submit.prevent="handleSubmit">
      <text-input v-model.trim="email" placeholder="Email" />
      <text-input
        v-model.trim="password"
        :type="showPassword ? 'text' : 'password'"
        placeholder="Password"
      />
      <input type="checkbox" v-model="showPassword" />
      <action-button text="Login" />
    </form>
    <router-link :to="{ name: 'SignUp' }"
      >New here? Click to create an account!</router-link
    >
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import TextInput from '@/components/TextInput.vue';
import ActionButton from '@/components/ActionButton.vue';
import axios, { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

interface User {
  email: string;
  username: string;
  mobileNumber: string;
}

const showPassword = ref(false);

const email = ref('');
const password = ref('');

const router = useRouter();
const userStore = useUserStore();

const handleSubmit = async () => {
  try {
    const body = {
      email: email.value,
      password: password.value,
    };

    const response = await axios.post(
      'http://localhost:8080/login',
      JSON.stringify(body),
      { withCredentials: true }
    );

    const user: User = {
      email: response?.data?.email,
      username: response?.data?.username,
      mobileNumber: '',
    };

    if (response.status === 200) {
      email.value = '';
      password.value = '';

      userStore.LOGIN_USER(user);

      // TODO: add user balance

      router.push({ name: 'Home' });
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else console.error('Unexpected error', error);
  }
};
</script>
