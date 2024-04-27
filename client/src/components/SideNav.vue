<template>
  <div class="mr-8 h-screen border-r-2 border-gray-200 dark:border-gray-500 px-12 pt-6">
    <ul class="flex cursor-pointer flex-col gap-6">
      <router-link :to="{ name: 'Home' }">
        <li class="grid grid-cols-3">
          <svg-icon type="mdi" :path="mdiHomeAccount" /><span>Home</span>
        </li>
      </router-link>
      <router-link :to="{ name: 'Balances' }">
        <li class="grid grid-cols-3">
          <svg-icon type="mdi" :path="mdiWallet" /><span>Balances</span>
        </li>
      </router-link>
      <router-link :to="{ name: 'Transactions' }">
        <li class="grid grid-cols-3">
          <svg-icon type="mdi" :path="mdiListBoxOutline" /><span
            >Transactions</span
          >
        </li></router-link
      >
      <router-link :to="{ name: 'Transfer' }">
        <li class="grid grid-cols-3">
          <svg-icon type="mdi" :path="mdiBankTransfer" /><span>Transfer</span>
        </li></router-link
      >
      <router-link :to="{ name: 'Wallets' }">
        <li class="grid grid-cols-3">
          <svg-icon type="mdi" :path="mdiCash" /><span>Wallets</span>
        </li>
      </router-link>
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
</script>
