<template>
  <div>
    <h1>Login Page</h1>
    <form @submit.prevent="handleSubmit">
      <label for="email">Email</label>
      <text-input id="email" v-model.trim="email" placeholder="Email" />
      <label for="password">Password</label>
      <text-input
        id="password"
        v-model.trim="password"
        placeholder="Password"
        :type="showPassword ? 'text' : 'password'"
      />
      <input type="checkbox" v-model="showPassword" />
      <action-button text="Login" />
    </form>
    <router-link :to="{ name: 'SignUp' }"
      >New here? Click to create an account!</router-link
    >
    <span>Temporary message: {{ responseMessage }}</span>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import TextInput from '@/components/TextInput.vue';
import ActionButton from '@/components/ActionButton.vue';
import { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import type { User } from '@/types/user';
import postLogin from '@/api/user';

const showPassword = ref(false);

const email = ref('');
const password = ref('');

const responseMessage = ref('');

const router = useRouter();
const userStore = useUserStore();

const handleSubmit = async () => {
  try {
    const body: { email: string; password: string } = {
      email: email.value,
      password: password.value,
    };

    const { data, status } = await postLogin(body);

    const user: User = {
      email: data?.email,
      username: data?.username,
      mobileNumber: '',
    };

    if (status === 200) {
      email.value = '';
      password.value = '';

      userStore.LOGIN_USER(user);

      // TODO: add user balance

      router.push({ name: 'Home' });
    }
  } catch (error: any) {
    if (error instanceof AxiosError) {
      if (error.response) {
        responseMessage.value = error.response?.data?.message;
      }
    } else {
      responseMessage.value = 'Unexpected error occurred';
    }
  }
};
</script>
