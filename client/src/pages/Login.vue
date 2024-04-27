<template>
  <div class="grid h-screen place-items-center">
    <div class="container w-1/3 rounded-lg border px-10 py-8 shadow-lg">
      <h1 class="mb-6 text-2xl font-bold">Login</h1>
      <form @submit.prevent="handleSubmit">
        <div class="mb-4 flex flex-col gap-2">
          <label for="email">Email</label>
          <text-input
            class="rounded-lg border px-4 py-2"
            id="email"
            v-model.trim="email"
            placeholder="leonlow@example.com"
          />
        </div>

        <div class="mb-6 flex flex-col gap-2">
          <label for="password">Password</label>
          <div class="relative flex w-full flex-wrap items-stretch">
            <text-input
              class="w-full rounded-lg border px-4 py-2"
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

        <action-button
          class="mb-4 w-full rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
          text="Login"
        />
      </form>

      <router-link
        :to="{ name: 'SignUp' }"
        class="flex justify-center text-cyan-900 underline underline-offset-4 transition hover:underline-offset-8"
        >New here? Click to create an account!</router-link
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { mdiEyeOutline, mdiEyeOffOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import { AxiosError } from "axios";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

import type { User, LOGIN_BODY } from "@/types/user";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { LOGIN } from "@/api/user";

// Data Fields
const showPassword = ref(false);

const email = ref("");
const password = ref("");

const responseMessage = ref("");

// Stores
const router = useRouter();
const userStore = useUserStore();

// Methods
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

const handleSubmit = async () => {
  try {
    const body: LOGIN_BODY = {
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
      email.value = "";
      password.value = "";

      userStore.LOGIN_USER(user);

      // TODO: add user balance

      alert("Logged In!");

      router.push({ name: "Home" });
    }
  } catch (error: any) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else {
      responseMessage.value = "Unexpected error occurred";
    }
  }
};
</script>
