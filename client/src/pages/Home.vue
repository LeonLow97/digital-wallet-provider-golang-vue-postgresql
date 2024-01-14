<template>
  <h1>Home Page</h1>
  <h3>Welcome Back {{ username }}</h3>
  <h3>Email {{ email }}</h3>
  <h3>Mobile Number {{ mobileNumber }}</h3>
  <action-button @click="handleClick" text="Logout" />
</template>

<script lang="ts" setup>
import axios from 'axios';
import { useRouter } from 'vue-router';
import ActionButton from '@/components/ActionButton.vue';
import { useUserStore } from '@/stores/user';

const router = useRouter();
const userStore = useUserStore();

const username = userStore.user.username;
const email = userStore.user.email;
const mobileNumber = userStore.user.mobileNumber;

const handleClick = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8080/logout',
      JSON.stringify({}),
      { withCredentials: true }
    );
    if (response.status === 200) {
      alert('logged out!');
      router.push({ name: 'Login' });
    }
  } catch (error) {
    alert(error);
  }
};
</script>
