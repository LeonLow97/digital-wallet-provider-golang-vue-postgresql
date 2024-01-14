<template>
  <h1>Home Page</h1>
  <action-button @click="handleClick" text="Logout" />
</template>

<script lang="ts" setup>
import axios from 'axios';
import { useRouter } from 'vue-router';
import ActionButton from '@/components/ActionButton.vue';

const router = useRouter();

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
