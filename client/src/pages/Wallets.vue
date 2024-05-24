<template>
  <h1 class="text-xl font-bold">Wallets</h1>
  <div class="grid grid-cols-3 gap-4">
    <span v-for="wallet of wallets" :key="wallet.walletID">
      <wallet-card
        :type="wallet.type"
        :amount="wallet.balance"
        :currency="wallet.currency"
      />
    </span>
  </div>
</template>
<script lang="ts" setup>
import WalletCard from "@/components/wallets/WalletCard.vue";
import { GET_WALLETS } from "@/api/wallet";
import { onMounted, ref } from "vue";
import type { Wallet } from "@/types/wallet";

// Change the type of wallets to match the expected structure
let wallets = ref([] as Wallet[]);

onMounted(async () => {
  const { data, status } = await GET_WALLETS();
  // Assuming data.wallets is an array of Wallet objects
  wallets.value = data.wallets;
});
</script>
