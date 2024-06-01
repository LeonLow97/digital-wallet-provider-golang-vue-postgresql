<template>
  <div class="flex items-center justify-between">
    <h1 class="text-xl font-bold tracking-wider">Balances</h1>
    <div class="flex gap-8">
      <action-button
        text="Deposit"
        class="rounded-lg bg-green-500 px-4 py-2 text-center text-white transition hover:bg-green-400"
        @click="handleDeposit()"
      />
      <action-button
        text="Withdraw"
        class="rounded-lg bg-pink-500 px-4 py-2 text-center text-white transition hover:bg-pink-400"
        @click="handleWithdraw()"
      />
      <action-button
        text="Currency Exchange"
        class="rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
        @click="handleCurrencyExchange()"
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

const userStore = useUserStore();
const currentAmount = ref(0);
let balances = ref([] as Balance[]);
let actionType = ref("");

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
  } catch (error: unknown) {
    alert(error);
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
