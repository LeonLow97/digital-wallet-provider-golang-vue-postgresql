<template>
  <h1>Home Page</h1>
  <h3>Welcome Back {{ username }}</h3>
  <h3>Email {{ email }}</h3>
  <h3>Mobile Number {{ mobileNumber }}</h3>
  <action-button @click="handleClick" text="Logout" />
</template>

<script lang="ts" setup>
import axios, { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import ActionButton from '@/components/ActionButton.vue';
import { useUserStore } from '@/stores/user';
import { onMounted, ref } from 'vue';

const router = useRouter();
const userStore = useUserStore();

const username = ref('');
const email = ref('');
const mobileNumber = ref('');

const handleClick = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8080/logout',
      JSON.stringify({}),
      { withCredentials: true }
    );
    if (response.status === 200) {
      userStore.LOGOUT_USER();

      alert('logged out!');
      router.push({ name: 'Login' });
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else console.error('Unexpected error', error);
  }
};

onMounted(() => {
  console.log('mounting', userStore.user.username);
  username.value = userStore.user.username;
  email.value = userStore.user.email;
  mobileNumber.value = userStore.user.mobileNumber;
});
</script>
