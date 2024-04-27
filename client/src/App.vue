<template>
  <main-nav v-if="isLoggedIn" />
  <div :class="dynamicClass">
    <side-nav class="col-span-2" v-if="isLoggedIn" />
    <router-view class="col-span-8" />
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/user";
import { computed } from "vue";
import MainNav from "@/components/MainNav.vue";
import SideNav from "@/components/SideNav.vue";

const userStore = useUserStore();

const isLoggedIn = computed(() => userStore.isLoggedIn);

const dynamicClass = computed(() => {
  return isLoggedIn.value ? "grid grid-cols-10" : null;
});
</script>
