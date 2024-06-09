<template>
  <div class="grid h-screen place-items-center">
    <div
      class="container w-1/3 rounded-lg border px-10 py-8 shadow-lg"
      v-if="!isLoggedIn"
    >
      <h1 class="mb-6 text-2xl font-bold">Login</h1>
      <form @submit.prevent="handleLogin">
        <div class="mb-4 flex flex-col gap-2">
          <label for="email">Email</label>
          <text-input
            id="email"
            v-model.trim="email"
            placeholder="leonlow@example.com"
          />
        </div>

        <div class="mb-6 flex flex-col gap-2">
          <label for="password">Password</label>
          <div class="relative flex w-full flex-wrap items-stretch">
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

        <router-link
          :to="{ name: 'ForgotPassword' }"
          class="text-cyan-900 underline underline-offset-4 dark:text-cyan-50"
          >Forgot your password?</router-link
        >

        <action-button
          class="mb-4 mt-4 w-full rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
          text="Login"
        />

        <router-link
          :to="{ name: 'SignUp' }"
          class="flex justify-center text-cyan-900 underline underline-offset-4 transition hover:underline-offset-8 dark:text-cyan-50"
          >New here? Click to create an account!</router-link
        >
      </form>
    </div>

    <div
      class="container w-3/5 rounded-lg border px-10 py-8 shadow-lg"
      v-if="showMFAConfigurationForm"
    >
      <configure-mfa
        :secret="mfaSecret"
        :url="mfaUrl"
        :email="email"
        @mfaConfigured="onMfaConfigured"
        v-if="showMFAConfigurationForm"
      />
    </div>

    <div
      class="container w-1/3 rounded-lg border px-10 py-8 shadow-lg"
      v-if="showMFAForm"
    >
      <verify-mfa :email="email" @mfaVerified="onMfaVerified" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { mdiEyeOutline, mdiEyeOffOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { useToastStore } from "@/stores/toast";

import type { User, LOGIN_REQUEST } from "@/types/user";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import VerifyMfa from "@/components/auth/VerifyMfa.vue";
import ConfigureMfa from "@/components/auth/ConfigureMfa.vue";
import { LOGIN } from "@/api/user";

// Data Fields
const showPassword = ref(false);

const email = ref("");
const password = ref("");
const isLoggedIn = ref(false);

const mfaUrl = ref("");
const mfaSecret = ref("");

const showMFAForm = ref(false);
const showMFAConfigurationForm = ref(false);

let user: User;

// Stores
const router = useRouter();
const userStore = useUserStore();
const toastStore = useToastStore();

// Methods
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

onMounted(async () => {
  userStore.LOGOUT_USER();
});

const handleLogin = async () => {
  try {
    const body: LOGIN_REQUEST = {
      email: email.value,
      password: password.value,
    };

    const { data, status } = await LOGIN(body);

    user = {
      firstName: data?.firstName,
      lastName: data?.lastName,
      email: data?.email,
      username: data?.username,
      sourceCurrency: data?.sourceCurrency,
      mobileCountryCode: data?.mobileCountryCode,
      mobileNumber: data?.mobileNumber,
    };

    if (status === 200) {
      isLoggedIn.value = true;

      if (data.isMfaConfigured) {
        showMFAForm.value = true;
      } else {
        showMFAConfigurationForm.value = true;
        mfaUrl.value = data.mfaConfig?.url!;
        mfaSecret.value = data.mfaConfig?.secret!;
      }
    }
  } catch (error: any) {
    if (error.code === "ERR_NETWORK") {
      toastStore.ERROR_TOAST("Network Error. Please try again later.");
      return;
    }
    toastStore.ERROR_TOAST(error.response?.data?.message, 2);
  }
};

const onMfaConfigured = (isConfigured: boolean) => {
  if (isConfigured) {
    showMFAConfigurationForm.value = false;
    // user is authenticated fully by password and mfa code
    userStore.LOGIN_USER(user);
    toastStore.SUCCESS_TOAST("Logged In Successfully!", 2);

    router.push({ name: "Home" });
  } else {
    router.push({ name: "Login" });
  }
};

const onMfaVerified = (isVerified: boolean) => {
  if (isVerified) {
    showMFAForm.value = false;
    // user is authenticated fully by password and mfa code
    userStore.LOGIN_USER(user);
    toastStore.SUCCESS_TOAST("Logged In Successfully!", 2);

    router.push({ name: "Home" });
  } else {
    router.push({ name: "Login" });
  }
};
</script>
