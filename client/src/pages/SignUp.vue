<template>
  <div class="grid h-screen place-items-center">
    <div class="container w-2/5 rounded-lg border px-10 py-8 shadow-lg">
      <h1 class="mb-8 text-2xl font-bold">Register an Account</h1>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-2 gap-x-8 gap-y-4">
          <text-input
            class="rounded-lg border px-4 py-2"
            v-model.trim="firstName"
            placeholder="First Name"
          />
          <text-input
            class="rounded-lg border px-4 py-2"
            v-model.trim="lastName"
            placeholder="Last Name"
          />
          <text-input
            class="col-span-2 rounded-lg border px-4 py-2"
            v-model="mobileNumber"
            placeholder="Mobile Number"
          />
          <text-input
            class="col-span-2 rounded-lg border px-4 py-2"
            v-model.trim="username"
            placeholder="Username"
          />
          <text-input
            class="col-span-2 rounded-lg border px-4 py-2"
            v-model="email"
            placeholder="Email"
          />
          <div class="relative col-span-2 flex w-full flex-wrap items-stretch">
            <text-input
              class="w-full rounded-lg border px-4 py-2"
              v-model="password"
              placeholder="Password"
              type="password"
            />
            <span
              class="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer"
            >
              <svg-icon
                v-model="showPassword"
                type="mdi"
                :path="showPassword ? mdiEyeOutline : mdiEyeOffOutline"
                @click="togglePasswordVisibility"
              ></svg-icon>
            </span>
          </div>
        </div>

        <div class="mt-5 flex items-center justify-between">
          <router-link
            :to="{ name: 'Login' }"
            class="text-cyan-900 underline underline-offset-4 transition hover:underline-offset-8"
            >Back to Login</router-link
          >
          <action-button
            text="Sign Up"
            class="rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { mdiEyeOutline, mdiEyeOffOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import { AxiosError } from "axios";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { SIGNUP } from "@/api/user";
import type { SIGNUP_BODY } from "@/types/user";

// Data Fields
const firstName = ref("");
const lastName = ref("");
const username = ref("");
const email = ref("");
const password = ref("");
const mobileNumber = ref("");

const showPassword = ref(false);

// Methods
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

const handleSubmit = async () => {
  try {
    const body: SIGNUP_BODY = {
      first_name: firstName.value === "" ? null : firstName.value,
      last_name: lastName.value === "" ? null : lastName.value,
      username: username.value,
      email: email.value,
      password: password.value,
      mobile_number: mobileNumber.value,
    };

    const { status } = await SIGNUP(body);

    if (status === 204) {
      firstName.value = "";
      lastName.value = "";
      username.value = "";
      email.value = "";
      password.value = "";
      mobileNumber.value = "";

      alert("Signed up successfully!");
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else console.error("Unexpected error", error);
  }
};
</script>
