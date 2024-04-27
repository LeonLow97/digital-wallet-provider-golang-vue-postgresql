<template>
  <div class="pt-6">
    <h3 class="text-center text-xl font-bold">
      Welcome Back <span class="uppercase">{{ username }}</span> !
    </h3>
  </div>
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
