<template>
  {{ isLoggedIn }}
  <main-nav v-if="isLoggedIn" />
  <router-view />
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user';
import { computed } from 'vue';
import MainNav from '@/components/MainNav.vue';

const userStore = useUserStore();

const isLoggedIn = computed(() => userStore.isLoggedIn);
</script>
