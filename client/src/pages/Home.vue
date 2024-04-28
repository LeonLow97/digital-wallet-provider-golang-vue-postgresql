<template>
  <h3 class="text-center text-2xl font-bold text-gray-900">
    Welcome Back <span class="uppercase text-blue-600">{{ username }}</span> !
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
