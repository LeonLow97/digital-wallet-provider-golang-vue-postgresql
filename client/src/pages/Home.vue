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
      'http://localhost:8080/api/v1/logout',
      JSON.stringify({}),
      { withCredentials: true }
    );
    if (response.status === 200) {
      userStore.LOGOUT_USER();
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        console.log(error.response);
        alert(error.response?.data?.message);
      }
    } else console.error('Unexpected error', error);
  } finally {
    // regardless of error, logout the user
    router.push({ name: 'Login' });
  }
};

onMounted(async () => {
  const response = await axios.get('http://localhost:8080/api/v1/wallet/all', {
    withCredentials: true,
  });

  console.log(response);

  username.value = userStore.user.username;
  email.value = userStore.user.email;
  mobileNumber.value = userStore.user.mobileNumber;
});
</script>
