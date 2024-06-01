<template>
  <div class="mt-2 flex items-center justify-between">
    <router-link
      :to="{ name: 'Wallets' }"
      class="tracking-wider text-blue-600 underline underline-offset-8 hover:text-blue-300"
      >&larr; Back to Wallets</router-link
    >
    <div>
      <action-button
        class="mr-8 rounded-lg border border-none bg-blue-500 px-8 py-2 text-center text-white transition hover:bg-blue-400"
        text="Make a Transfer"
        @click="handleTransfer"
      />
      <action-button
        class="mr-8 rounded-lg border border-none bg-green-500 px-8 py-2 text-center text-white transition hover:bg-green-400"
        text="Top Up Wallet"
        @click="handleWalletExchanges('Top Up')"
      />
      <action-button
        class="rounded-lg border border-none bg-pink-600 px-8 py-2 text-center text-white transition hover:bg-pink-500"
        text="Cash Out to Main Balance"
        @click="handleWalletExchanges('Cash Out')"
      />
    </div>
  </div>

  <!-- Wallet Card -->
  <div class="w-full rounded-lg border border-gray-200 bg-white p-6 shadow-md">
    <h2 class="mb-4 text-lg font-bold capitalize text-blue-800">
      {{ wallet?.walletType }}
    </h2>
    <h3 class="text-md mb-2 font-bold capitalize text-gray-500">
      {{ wallet?.currencyAmount?.length }} Available Currencies:
    </h3>
    <div
      class="grid w-1/6 grid-cols-2 gap-4"
      v-for="item in wallet?.currencyAmount"
    >
      <p class="text-gray-600">{{ item.amount }}</p>
      <p class="text-blue-600">{{ item.currency }}</p>
    </div>
  </div>

  <wallet-exchanges-modal
    @close-modal="closeModal"
    @form-submitted="formSubmitted"
    :open-modal="openModal"
    :wallet-id="walletId"
    :wallet-currencies="walletCurrencies!"
    :action-type="actionType"
  />

  <create-transaction-modal
    @close-transaction-modal="closeTransactionModal"
    @form-submitted="formSubmitted"
    :open-transaction-modal="openTransactionModal"
    :wallet-currencies="walletCurrencies!"
    :wallet-id="walletId!"
  />
</template>

<script lang="ts" setup>
import { useRoute, useRouter } from "vue-router";
import { onMounted, ref } from "vue";
import ActionButton from "@/components/ActionButton.vue";
import { GET_WALLET } from "@/api/wallet";
import type { Wallet } from "@/types/wallet";
import WalletExchangesModal from "@/components/wallets/WalletExchangesModal.vue";
import CreateTransactionModal from "@/components/transactions/CreateTransactionModal.vue";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const route = useRoute();
const router = useRouter();

const wallet = ref<Wallet | null>(null);
const walletId = ref<number | null>(null);
const walletCurrencies = ref<string[]>([]);
const actionType = ref("");

const openModal = ref(false);
const openTransactionModal = ref(false);

const closeModal = () => {
  openModal.value = false;
};

const closeTransactionModal = () => {
  openTransactionModal.value = false;
};

const formSubmitted = (val: string) => {
  getWallet();

  // clean data
  walletCurrencies.value = [];

  const msg = `${val} Successfully!`;
  toastStore.SUCCESS_TOAST(msg);
};

onMounted(() => {
  getWallet();
});

const getWallet = async () => {
  // ensure that params is a number
  const id = Number(route.params.id);
  if (isNaN(id)) {
    router.push({ name: "Wallets" });
    return;
  }
  walletId.value = id;

  try {
    const { data, status } = await GET_WALLET(id);

    if (status === 200) {
      wallet.value = data;
      wallet.value.currencyAmount.forEach((item) => {
        walletCurrencies.value?.push(item.currency);
      });
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const handleWalletExchanges = (action: string) => {
  openModal.value = true;
  actionType.value = action;
};

const handleTransfer = () => {
  openTransactionModal.value = true;
};
</script>
