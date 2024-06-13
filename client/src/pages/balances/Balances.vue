<template>
  <div class="flex items-center justify-between">
    <h1 class="text-xl font-bold tracking-wider">Balances</h1>
    <div class="flex gap-8">
      <action-button
        text="Currency Exchange"
        class="rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-300 dark:hover:bg-blue-700"
        @click="handleCurrencyExchange()"
      />
      <action-button
        text="Deposit"
        class="rounded-lg bg-slate-700 px-4 py-2 text-center text-white transition hover:bg-slate-500 dark:bg-slate-600 dark:hover:bg-slate-700"
        @click="handleDeposit()"
      />
      <action-button
        text="Withdraw"
        class="rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300"
        @click="handleWithdraw()"
      />
    </div>
  </div>
  <balances-table @form-submitted="formSubmitted" :balances="balances" />

  <balance-modal
    @close-modal="closeModal"
    @form-submitted="formSubmitted"
    class="text-left"
    :open-modal="openModal"
    :action-type="actionType"
    :currency="userStore.user.sourceCurrency"
    :current-amount="currentAmount"
  />

  <currency-exchange-modal
    @close-modal="closeCurrencyExchangeModal"
    @form-submitted="formSubmitted"
    :open-modal="openCurrencyExchangeModal"
    :currency="userStore.user.sourceCurrency"
  />
</template>

<script lang="ts" setup>
import { GET_BALANCES } from "@/api/balances";
import type { Balance } from "@/types/balances";
import { onMounted, ref } from "vue";
import { useUserStore } from "@/stores/user";

import ActionButton from "@/components/ActionButton.vue";
import BalancesTable from "@/components/balances/BalancesTable.vue";
import BalanceModal from "@/components/balances/BalanceModal.vue";
import CurrencyExchangeModal from "@/components/balances/CurrencyExchangeModal.vue";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const userStore = useUserStore();
const currentAmount = ref(0);
const balances = ref([] as Balance[]);
const actionType = ref("");

const openModal = ref(false);
const openCurrencyExchangeModal = ref(false);

onMounted(async () => {
  getBalances();
});

const formSubmitted = () => {
  getBalances();
};

const getBalances = async () => {
  try {
    const { data, status } = await GET_BALANCES();

    // sort the balances by amount in descending order
    balances.value = data.balances.sort((a, b) => b.balance - a.balance);

    balances.value.forEach((balance) => {
      if (balance.currency === userStore.user.sourceCurrency) {
        currentAmount.value = balance.balance;
      }
    });
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const closeModal = () => {
  openModal.value = false;
};

const closeCurrencyExchangeModal = () => {
  openCurrencyExchangeModal.value = false;
};

const handleDeposit = () => {
  actionType.value = "deposit";
  openModal.value = true;
};

const handleWithdraw = () => {
  actionType.value = "withdraw";
  openModal.value = true;
};

const handleCurrencyExchange = async () => {
  openCurrencyExchangeModal.value = true;
};
</script>
