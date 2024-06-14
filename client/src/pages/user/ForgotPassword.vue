<template>
  <div class="grid h-screen place-items-center">
    <div class="container w-2/5 rounded-lg border px-10 py-8 shadow-lg">
      <h1 class="mb-8 text-2xl font-bold">Forgot Password</h1>
      <form @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-6">
          <text-input type="email" v-model.trim="email" placeholder="Email" />

          <action-button
            text="Submit"
            class="w-full rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
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
import { ref } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { SEND_PASSWORD_RESET_EMAIL } from "@/api/user";
import type { SendPasswordResetEmailRequest } from "@/types/user";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();

const email = ref("");

const handleSubmit = async () => {
  try {
    const body: SendPasswordResetEmailRequest = {
      email: email.value,
    };

    const { status } = await SEND_PASSWORD_RESET_EMAIL(body);

    if (status === 204) {
      email.value = "";

      toastStore.SUCCESS_TOAST(
        "Password reset email sent successfully! If the email address is registered, you will receive an email with instructions to reset your password.",
      );
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
