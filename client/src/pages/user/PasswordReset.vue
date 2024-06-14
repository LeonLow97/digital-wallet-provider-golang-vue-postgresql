<template>
  <div class="grid h-screen place-items-center">
    <div class="container w-2/5 rounded-lg border px-10 py-8 shadow-lg">
      <h1 class="mb-8 text-2xl font-bold">Password Reset</h1>
      <form @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-6">
          <text-input
            type="password"
            v-model.trim="password"
            placeholder="New Password"
          />

          <action-button
            text="Submit"
            class="w-full rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
          />
        </div>
      </form>

      <router-link
        :to="{ name: 'Login' }"
        class="mt-4 flex justify-center text-cyan-900 underline underline-offset-4 transition hover:underline-offset-8 dark:text-cyan-50"
        >&larr; Back to Login</router-link
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from "vue-router";
import { ref, onMounted } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { PASSWORD_RESET } from "@/api/user";
import type { PasswordResetRequest } from "@/types/user";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const route = useRoute();

const token = ref("");
const password = ref("");

onMounted(async () => {
  const tokenParam = route.params.token;

  if (typeof tokenParam === "string") {
    token.value = decodeURIComponent(tokenParam);
  } else if (Array.isArray(tokenParam)) {
    token.value = decodeURIComponent(tokenParam[0]);
  }
});

const handleSubmit = async () => {
  try {
    const body: PasswordResetRequest = {
      token: token.value,
      password: password.value,
    };

    const { status } = await PASSWORD_RESET(body);

    if (status === 204) {
      password.value = "";

      toastStore.SUCCESS_TOAST("Successfully Reset Password!");
    }
  } catch (error: any) {
    if (error.code === "ERR_NETWORK") {
      toastStore.ERROR_TOAST("Network Error. Please try again later.");
      return;
    }
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
