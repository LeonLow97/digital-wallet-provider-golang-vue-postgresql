<template>
  <h1>Sign Up Page</h1>
  <form @submit.prevent="handleSubmit">
    <text-input v-model.trim="firstName" placeholder="First Name" />
    <text-input v-model.trim="lastName" placeholder="Last Name" />
    <text-input v-model.trim="username" placeholder="Username" />
    <text-input v-model="email" placeholder="Email" />
    <text-input v-model="password" placeholder="Password" type="password" />
    <text-input v-model="mobileNumber" placeholder="Mobile Number" />

    <action-button text="Sign Up" />
  </form>
  <router-link :to="{ name: 'Login' }">Back to Login</router-link>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import TextInput from '@/components/TextInput.vue';
import ActionButton from '@/components/ActionButton.vue';
import axios from 'axios';

const firstName = ref('');
const lastName = ref('');
const username = ref('');
const email = ref('');
const password = ref('');
const mobileNumber = ref('');

const handleSubmit = async () => {
  try {
    const body = {
      first_name: firstName.value === '' ? null : firstName.value,
      last_name: lastName.value === '' ? null : lastName.value,
      username: username.value,
      email: email.value,
      password: password.value,
      mobile_number: mobileNumber.value,
    };

    const response = await axios.post(
      'http://localhost:8080/signup',
      JSON.stringify(body),
      { withCredentials: true }
    );

    if (response.status) {
      firstName.value = '';
      lastName.value = '';
      username.value = '';
      email.value = '';
      password.value = '';
      mobileNumber.value = '';

      alert('sign up!');
    }
  } catch (error: any) {
    alert(error.response?.data?.message || 'An unexpected error occurred');
  }
};
</script>
