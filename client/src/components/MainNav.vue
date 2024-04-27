<template>
  <div class="flex justify-between bg-blue-900 p-4 text-white">
    <h2 class="pl-8 font-mono text-lg font-bold uppercase">Mobile Wallet</h2>
    <div class="pr-8">
      <ul class="flex gap-x-8">
        <li class="cursor-pointer hover:bg-blue-500">
          <svg-icon type="mdi" :path="mdiAccount" />
        </li>
        <li class="cursor-pointer hover:bg-blue-500" @click="handleLogout">
          <svg-icon type="mdi" :path="mdiLogout" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { mdiLogout, mdiAccount } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
import { AxiosError } from "axios";
import { LOGOUT } from "@/api/user";
import { useUserStore } from "@/stores/user";
import { useRouter } from "vue-router";

const router = useRouter();
const userStore = useUserStore();

const handleLogout = async () => {
  try {
    const { status } = await LOGOUT();

    if (status !== 200) {
      console.log("logout was unsuccessful");
    }

    alert("logged out!");
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        console.log(error.response);
        alert(error.response?.data?.message);
      }
    } else console.error("Unexpected error", error);
  } finally {
    // regardless of error, logout the user
    userStore.LOGOUT_USER();
    router.push({ name: "Login" });
  }
};
</script>
