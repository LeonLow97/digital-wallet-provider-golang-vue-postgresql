<template>
  <div class="login_container">
    <h1>Login</h1>
    <form @submit.prevent="handleSubmit">
      <div>
        <label for="email">Email</label>
        <text-input
          class="text_input"
          id="email"
          v-model.trim="email"
          placeholder="leonlow@example.com"
        />
      </div>

      <div>
        <label for="password">Password</label>
        <div class="text_input-container">
          <text-input
            class="text_input"
            id="password"
            v-model.trim="password"
            placeholder="Password"
            :type="showPassword ? 'text' : 'password'"
          />
          <span class="icon">
            <svg-icon
              v-model="showPassword"
              type="mdi"
              :path="showPassword ? mdiEyeOutline : mdiEyeOffOutline"
              @click="togglePasswordVisibility"
            ></svg-icon>
          </span>
        </div>
      </div>

      <div class="button_container">
        <action-button text="Login" />
      </div>
    </form>
    <router-link :to="{ name: 'SignUp' }"
      >New here? Click to create an account!</router-link
    >
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { mdiEyeOutline, mdiEyeOffOutline } from '@mdi/js';
import SvgIcon from '@jamescoyle/vue-icon';
import { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

import type { User } from '@/types/user';
import TextInput from '@/components/TextInput.vue';
import ActionButton from '@/components/ActionButton.vue';
import LOGIN from '@/api/user';

const showPassword = ref(false);

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

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

    const { data, status } = await LOGIN(body);

    const user: User = {
      email: data?.email,
      username: data?.username,
      mobileNumber: data?.mobileNumber,
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

<style scoped>
.login_container {
  display: flex;
  flex-direction: column;
  gap: 2.4rem;
  padding: 4.8rem;
  border-radius: 2.4rem;
  box-shadow: 0 1.2rem 2.4rem rgba(0, 0, 0, 0.2);
  width: 50rem;

  /** Center login container */
  margin: 0 auto;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.text_input {
  border-top: none;
  border-left: none;
  border-right: none;
  padding: 1.2rem 1rem 1.2rem 1rem;
  font-size: 1.8rem;
}

.text_input-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: row;
}

.text_input-container .text_input {
  width: 90%;
}

.button_container {
  display: flex;
  align-items: center;
  justify-content: center;
}

form {
  display: flex;
  flex-direction: column;
  gap: 2.8rem;
}

form div {
  display: flex;
  flex-direction: column;
}

.icon {
  cursor: pointer;
}
</style>
