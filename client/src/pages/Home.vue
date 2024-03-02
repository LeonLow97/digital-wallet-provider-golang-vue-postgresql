<template>
  <h3>Welcome Back {{ username }}</h3>
  <h3>Email {{ email }}</h3>
  <h3>Mobile Number {{ mobileNumber }}</h3>
  <action-button @click="handleLogout" text="Logout" />
</template>

<script lang="ts" setup>
import { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import ActionButton from '@/components/ActionButton.vue';
import { useUserStore } from '@/stores/user';
import { onMounted, ref } from 'vue';
import { LOGOUT } from '@/api/user';

const router = useRouter();
const userStore = useUserStore();

const username = ref('');
const email = ref('');
const mobileNumber = ref('');

const handleLogout = async () => {
  try {
    const { status } = await LOGOUT();

    if (status !== 200) {
      console.log('logout was unsuccessful');
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
    userStore.LOGOUT_USER();
    router.push({ name: 'Login' });
  }
};

onMounted(async () => {
  username.value = userStore.user.username;
  email.value = userStore.user.email;
  mobileNumber.value = userStore.user.mobileNumber;
});
</script>
