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
import type { SEND_PASSWORD_RESET_EMAIL_REQUEST } from "@/types/user";
import { AxiosError } from "axios";

const email = ref("");
const responseMessage = ref("");

const handleSubmit = async () => {
  try {
    const body: SEND_PASSWORD_RESET_EMAIL_REQUEST = {
      email: email.value,
    };

    const { status } = await SEND_PASSWORD_RESET_EMAIL(body);

    if (status === 204) {
      email.value = "";

      alert(
        "Password reset email sent successfully! If the email address is registered, you will receive an email with instructions to reset your password.",
      );
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
