<template>
  <main-nav @toggle-side-nav="toggleSideNav" v-if="isLoggedIn" />
  <div :class="dynamicClass" class="relative">
    <side-nav
      :class="sideNavDynamicCss"
      v-if="isLoggedIn"
      :openSideNav="openSideNav!"
    />
    <div
      :class="routerViewDynamicCss"
      class="flex h-screen flex-col gap-4 pr-8 pt-6"
    >
      <router-view />
    </div>

    <toast
      v-if="toast.visible"
      :message="toast.message"
      :toastType="toast.toastType"
    />
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/user";
import { useToastStore } from "@/stores/toast";
import { ref, computed } from "vue";
import MainNav from "@/components/MainNav.vue";
import SideNav from "@/components/SideNav.vue";
import Toast from "@/components/Toast.vue";

const userStore = useUserStore();
const toast = useToastStore().toast;

const sideNavDynamicCss = ref("col-span-2");
const routerViewDynamicCss = ref("col-span-8");
const openSideNav = ref<boolean>(true);

const isLoggedIn = computed(() => userStore.isLoggedIn);

const dynamicClass = computed(() => {
  return isLoggedIn.value ? "grid grid-cols-10" : null;
});

const toggleSideNav = (toggle: boolean) => {
  openSideNav.value = toggle;
  if (toggle) {
    sideNavDynamicCss.value = "col-span-2";
    routerViewDynamicCss.value = "col-span-8";
  } else {
    sideNavDynamicCss.value = "col-span-1";
    routerViewDynamicCss.value = "col-span-9";
  }
};
</script>
