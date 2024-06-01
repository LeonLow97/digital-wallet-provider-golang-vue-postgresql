<template>
  <div class="grid h-screen place-items-center">
    <div class="container w-2/5 rounded-lg border px-10 py-8 shadow-lg">
      <h1 class="mb-8 text-2xl font-bold">Register an Account</h1>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-2 gap-x-8 gap-y-6">
          <text-input v-model.trim="firstName" placeholder="First Name" />
          <text-input v-model.trim="lastName" placeholder="Last Name" />

          <div class="col-span-2 flex gap-4">
            <select
              class="rounded-md border border-gray-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50"
              v-model.trim="mobileCountryCode"
            >
              <option value="+65">+ 65</option>
              <option value="+60">+ 60</option>
              <option value="+61">+ 61</option>
              <option value="+1">+ 1</option>
            </select>

            <text-input
              class="col-span-2"
              v-model="mobileNumber"
              placeholder="Mobile Number"
            />
          </div>

          <text-input
            class="col-span-2"
            v-model.trim="username"
            placeholder="Username"
          />
          <text-input class="col-span-2" v-model="email" placeholder="Email" />
          <div class="relative col-span-2 flex w-full flex-wrap items-stretch">
            <text-input
              id="password"
              v-model.trim="password"
              placeholder="Password"
              :type="showPassword ? 'text' : 'password'"
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
            class="text-cyan-900 underline underline-offset-4 transition hover:underline-offset-8 dark:text-cyan-50"
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
import type { SIGNUP_REQUEST } from "@/types/user";

// Data Fields
const firstName = ref("");
const lastName = ref("");
const username = ref("");
const email = ref("");
const password = ref("");
const mobileCountryCode = ref("+65");
const mobileNumber = ref("");

const showPassword = ref(false);

// Methods
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

const handleSubmit = async () => {
  try {
    const body: SIGNUP_REQUEST = {
      first_name: firstName.value === "" ? null : firstName.value,
      last_name: lastName.value === "" ? null : lastName.value,
      username: username.value,
      email: email.value,
      password: password.value,
      mobile_country_code: mobileCountryCode.value,
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
  } catch (error: any) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else console.error("Unexpected error", error);
  }
};
</script>
