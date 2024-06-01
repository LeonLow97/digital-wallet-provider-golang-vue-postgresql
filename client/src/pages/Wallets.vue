<template>
  <h1 class="text-xl font-bold tracking-wider">Wallets</h1>
  <div class="mb-4">
    <action-button
      class="rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
      text="Create Wallet"
      @click="handleCreateWallet"
    />
  </div>

  <div class="grid auto-rows-fr grid-cols-3 gap-8">
    <div v-for="wallet of wallets" :key="wallet.id" class="flex">
      <wallet-card
        @click="handleClickViewWallet(wallet.id)"
        :type="wallet.walletType"
        :currency-amount="wallet.currencyAmount"
        class="h-full flex-grow transition duration-300 ease-in-out hover:scale-y-105"
      />
    </div>
  </div>

  <create-wallet-modal
    @close-modal="closeModal"
    @form-submitted="formSubmitted"
    :open-modal="openModal"
  />
</template>
<script lang="ts" setup>
import { useRouter } from "vue-router";
import WalletCard from "@/components/wallets/WalletCard.vue";
import { GET_WALLETS } from "@/api/wallet";
import { onMounted, ref } from "vue";
import type { Wallet } from "@/types/wallet";
import ActionButton from "@/components/ActionButton.vue";
import CreateWalletModal from "@/components/wallets/CreateWalletModal.vue";

const router = useRouter();
const openModal = ref(false);

// Change the type of wallets to match the expected structure
let wallets = ref([] as Wallet[]);

onMounted(async () => {
  getWallets();
});

const getWallets = async () => {
  try {
    const { data, status } = await GET_WALLETS();

    wallets.value = data;
  } catch (error: any) {
    alert(error?.response.data.message);
  }
};

const handleClickViewWallet = (walletId: number) => {
  router.push({ name: "Wallet", params: { id: walletId } });
};

const handleCreateWallet = () => {
  openModal.value = true;
};

const closeModal = () => {
  openModal.value = false;
};

const formSubmitted = () => {
  getWallets();
};
</script>
