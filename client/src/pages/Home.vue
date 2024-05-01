<template>
  <h3 class="text-xl font-bold text-gray-900 dark:text-white">
    Welcome Back, <br />
    <span class="text-2xl uppercase text-blue-600 dark:text-blue-500">{{
      username
    }}</span>
  </h3>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/stores/user";
import { onMounted, ref } from "vue";

const userStore = useUserStore();

const username = ref("");
const email = ref("");
const mobileNumber = ref("");

onMounted(async () => {
  username.value = userStore.user.username;
  email.value = userStore.user.email;
  mobileNumber.value = userStore.user.mobileNumber;
});
</script>
