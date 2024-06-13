<template>
  <div class="flex justify-between bg-slate-800 text-white">
    <div class="flex items-center gap-8 pl-12">
      <svg-icon
        type="mdi"
        :path="openSideNav ? mdiMenuOpen : mdiMenuClose"
        class="cursor-pointer"
        @click="toggleSideNav"
      />
      <router-link
        class="cursor-pointer font-mono text-lg font-bold uppercase tracking-wider"
        :to="{ name: 'Home' }"
      >
        Digital Wallet
      </router-link>
    </div>
    <div class="pr-8">
      <ul class="flex">
        <router-link
          :to="{ name: 'UserProfile' }"
          active-class="active"
          class="nav-link w-full"
        >
          <li class="flex w-full px-6 py-5 transition hover:bg-blue-500">
            <svg-icon type="mdi" :path="mdiAccount" />
          </li>
        </router-link>
        <router-link
          :to="{ name: 'Settings' }"
          active-class="active"
          class="nav-link w-full"
        >
          <li class="flex w-full px-6 py-5 transition hover:bg-blue-500">
            <svg-icon type="mdi" :path="mdiCogOutline" />
          </li>
        </router-link>
        <li
          class="flex w-full cursor-pointer px-6 py-5 transition hover:bg-blue-500"
          @click="handleLogout"
        >
          <svg-icon type="mdi" :path="mdiLogout" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  mdiAccount,
  mdiCogOutline,
  mdiLogout,
  mdiMenuOpen,
  mdiMenuClose,
} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import { LOGOUT } from "@/api/user";
import { useUserStore } from "@/stores/user";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useToastStore } from "@/stores/toast";

const openSideNav = ref(true);

const router = useRouter();
const userStore = useUserStore();

const emits = defineEmits(["toggleSideNav"]);

const toastStore = useToastStore();

const handleLogout = async () => {
  try {
    const { status } = await LOGOUT();

    if (status === 200) {
      toastStore.SUCCESS_TOAST("Logged Out! We hope to see you again!");
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST("Internal Server Error", 2);
  } finally {
    // regardless of error, logout the user
    userStore.LOGOUT_USER();
    router.push({ name: "Login" });
  }
};

const toggleSideNav = () => {
  openSideNav.value = !openSideNav.value;
  emits("toggleSideNav", openSideNav.value);
};
</script>

<style scoped>
.nav-link.active {
  background-color: #1d4ed8; /* Tailwind's bg-blue-900 */
  color: white;
}
</style>
