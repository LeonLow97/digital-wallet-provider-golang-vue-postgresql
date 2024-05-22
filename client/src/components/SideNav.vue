<template>
  <div
    class="mr-8 h-screen border-r-2 border-gray-200 py-2 dark:border-gray-500"
  >
    <ul class="flex cursor-pointer flex-col">
      <router-link :to="{ name: 'Home' }">
        <li
          class="grid grid-cols-3 px-12 py-4 transition hover:bg-blue-500 hover:text-white"
        >
          <svg-icon type="mdi" :path="mdiHomeAccount" /><span
            :class="props.openSideNav ? '' : 'hidden'"
            >Home</span
          >
        </li>
      </router-link>
      <router-link :to="{ name: 'Balances' }">
        <li
          class="grid grid-cols-3 px-12 py-4 transition hover:bg-blue-500 hover:text-white"
        >
          <svg-icon type="mdi" :path="mdiWallet" /><span
            :class="props.openSideNav ? '' : 'hidden'"
            >Balances</span
          >
        </li>
      </router-link>
      <router-link :to="{ name: 'Wallets' }">
        <li
          class="grid grid-cols-3 px-12 py-4 transition hover:bg-blue-500 hover:text-white"
        >
          <svg-icon type="mdi" :path="mdiCash" /><span
            :class="props.openSideNav ? '' : 'hidden'"
            >Wallets</span
          >
        </li>
      </router-link>
      <router-link :to="{ name: 'Transactions' }">
        <li
          class="grid grid-cols-3 px-12 py-4 transition hover:bg-blue-500 hover:text-white"
        >
          <svg-icon type="mdi" :path="mdiListBoxOutline" /><span
            :class="props.openSideNav ? '' : 'hidden'"
            >Transactions</span
          >
        </li></router-link
      >
    </ul>
  </div>
</template>

<script lang="ts" setup>
import {
  mdiHomeAccount,
  mdiWallet,
  mdiListBoxOutline,
  mdiBankTransfer,
  mdiCash,
} from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";

const props = defineProps<{
  openSideNav: boolean;
}>();
</script>
