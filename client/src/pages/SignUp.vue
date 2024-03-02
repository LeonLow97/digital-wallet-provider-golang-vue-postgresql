<template>
  <div class="signup_container">
    <h1>Register an Account</h1>
    <form @submit.prevent="handleSubmit">
      <div>
        <text-input
          class="text_input"
          v-model.trim="firstName"
          placeholder="First Name"
        />
        <text-input
          class="text_input"
          v-model.trim="lastName"
          placeholder="Last Name"
        />
        <text-input
          class="text_input"
          v-model.trim="username"
          placeholder="Username"
        />
        <text-input
          class="text_input"
          v-model="mobileNumber"
          placeholder="Mobile Number"
        />
        <text-input class="text_input" v-model="email" placeholder="Email" />
        <text-input
          class="text_input"
          v-model="password"
          placeholder="Password"
          type="password"
        />
      </div>

      <div class="button_container">
        <action-button text="Sign Up" />
      </div>
    </form>
    <router-link :to="{ name: 'Login' }">Back to Login</router-link>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { AxiosError } from 'axios';
import TextInput from '@/components/TextInput.vue';
import ActionButton from '@/components/ActionButton.vue';
import { SIGNUP } from '@/api/user';
import type { SIGNUP_BODY } from '@/types/user';

const firstName = ref('');
const lastName = ref('');
const username = ref('');
const email = ref('');
const password = ref('');
const mobileNumber = ref('');

const handleSubmit = async () => {
  try {
    const body: SIGNUP_BODY = {
      first_name: firstName.value === '' ? null : firstName.value,
      last_name: lastName.value === '' ? null : lastName.value,
      username: username.value,
      email: email.value,
      password: password.value,
      mobile_number: mobileNumber.value,
    };

    const { status } = await SIGNUP(body);

    if (status === 204) {
      firstName.value = '';
      lastName.value = '';
      username.value = '';
      email.value = '';
      password.value = '';
      mobileNumber.value = '';

      alert('Signed up successfully!');
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

<style scoped>
.signup_container {
  display: flex;
  flex-direction: column;
  gap: 2.4rem;
  width: 80rem;

  padding: 4.8rem;
  border-radius: 2.4rem;
  box-shadow: 0 1.2rem 2.4rem rgba(0, 0, 0, 0.2);

  /** Center signup container */
  margin: 0 auto;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

form {
  display: flex;
  flex-direction: column;
  gap: 3.2rem;
}

form div {
  display: grid;
  grid-template-columns: 1fr 1fr;
  row-gap: 3.2rem;
  column-gap: 2.4rem;
}

.text_input {
  border-top: none;
  border-left: none;
  border-right: none;
  padding: 1.2rem 1rem 1.2rem 1rem;
  font-size: 1.8rem;
}

.button_container {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
