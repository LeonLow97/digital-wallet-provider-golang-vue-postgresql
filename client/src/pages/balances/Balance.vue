<template>
  <div class="mx-auto w-full">
    <div class="mt-2 flex items-center justify-between">
      <router-link
        :to="{ name: 'Balances' }"
        class="tracking-wider text-blue-600 underline underline-offset-8 hover:text-blue-300"
        >&larr; Back to Balances</router-link
      >
      <div>
        <action-button
          class="mr-8 rounded-lg border border-none bg-green-500 px-8 py-2 text-center text-white transition hover:bg-green-400"
          text="Deposit"
          @click="handleDeposit"
        />
        <action-button
          class="rounded-lg border border-none bg-pink-600 px-8 py-2 text-center text-white transition hover:bg-pink-400"
          text="Withdraw"
          @click="handleWithdraw"
        />
      </div>
    </div>

    <!-- Balance Card -->
    <div class="mt-8 flex items-start justify-start gap-4">
      <div
        class="w-full rounded-lg border border-gray-200 bg-white p-6 shadow-md"
      >
        <h2 class="mb-4 text-lg font-bold text-gray-800">Balance</h2>
        <div class="grid w-2/5 grid-cols-2 gap-4">
          <p class="text-gray-600">Balance:</p>
          <p class="text-blue-600">{{ balance?.balance }}</p>

          <p class="text-gray-600">Currency:</p>
          <p class="text-blue-600">{{ balance?.currency }}</p>

          <p class="text-gray-600">Created At:</p>
          <p class="italic text-gray-600">
            {{ formatDate(balance?.createdAt) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Balance History Table -->
    <h2 class="mb-4 mt-8 text-lg font-bold">Balance History</h2>
    <div class="mb-14 shadow-md">
      <table class="w-full table-fixed text-sm">
        <thead>
          <tr class="bg-gray-100">
            <th class="px-4 py-2">Amount</th>
            <th class="px-4 py-2">Currency</th>
            <th class="px-4 py-2">Type</th>
            <th class="px-4 py-2">Created At</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="history in balanceHistory?.balanceHistory"
            class="text-center"
          >
            <td
              class="px-4 py-2"
              :class="
                history.type.trim().toLowerCase() === 'deposit'
                  ? 'text-green-600'
                  : 'text-red-600'
              "
            >
              {{ history.amount }}
            </td>
            <td class="px-4 py-2">{{ history.currency }}</td>
            <td
              class="px-4 py-2 font-bold uppercase"
              :class="
                history.type.trim().toLowerCase() === 'deposit'
                  ? 'text-green-600'
                  : 'text-red-600'
              "
            >
              {{ history.type }}
            </td>
            <td class="px-4 py-2">{{ formatDate(history.createdAt) }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="!balanceHistory?.balanceHistory" class="p-4 text-center">
        No Balance History Found!
      </div>
    </div>
  </div>

  <balance-modal
    @close-modal="closeModal"
    @form-submitted="formSubmitted"
    :open-modal="openModal"
    :action-type="actionType"
    :currency="mainBalance.currency"
    :current-amount="mainBalance.amount"
  />
</template>

<script lang="ts" setup>
import { useUserStore } from "@/stores/user";
import { onMounted, ref } from "vue";
import { GET_BALANCE, GET_BALANCE_HISTORY, GET_BALANCES } from "@/api/balances";
import { useRoute, useRouter } from "vue-router";
import type { Balance, GetBalanceHistoryResponse } from "@/types/balances";
import { format } from "date-fns";
import ActionButton from "@/components/ActionButton.vue";
import BalanceModal from "@/components/balances/BalanceModal.vue";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const mainBalance = ref<{ amount: number; currency: string }>({
  amount: 0,
  currency: "",
});
const balance = ref<Balance>({
  id: 0,
  balance: 0,
  currency: "",
  createdAt: "",
  updatedAt: "",
});
const balanceHistory = ref<GetBalanceHistoryResponse | null>(null);
let actionType = ref("");
const openModal = ref(false);

onMounted(() => {
  getBalanceAndBalanceHistory();
  getBalances();
});

const formSubmitted = () => {
  getBalanceAndBalanceHistory();
};

const getBalanceAndBalanceHistory = async () => {
  // ensure that params is a number
  const id = Number(route.params.id);
  if (isNaN(id)) {
    router.push({ name: "Balances" });
    return;
  }

  try {
    // calling these 2 endpoint asynchronously since they are independent
    const [balanceResponse, balanceHistoryResponse] = await Promise.all([
      GET_BALANCE(id),
      GET_BALANCE_HISTORY(id),
    ]);

    if (balanceResponse.status === 200) {
      balance.value = balanceResponse.data;
    }

    if (balanceHistoryResponse.status === 200) {
      balanceHistory.value = balanceHistoryResponse.data;
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const getBalances = async () => {
  try {
    const { data, status } = await GET_BALANCES();

    let filteredMainBalance = data.balances.filter(
      (item) => item.currency === userStore.user.sourceCurrency,
    );
    mainBalance.value = {
      amount: filteredMainBalance[0].balance,
      currency: filteredMainBalance[0].currency,
    };
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const formatDate = (dateString: string | undefined): string => {
  if (dateString) {
    return format(new Date(dateString), "PPpp");
  } else {
    return "Invalid Date Format";
  }
};

const handleDeposit = () => {
  actionType.value = "deposit";
  openModal.value = true;
};

const handleWithdraw = () => {
  actionType.value = "withdraw";
  openModal.value = true;
};

const closeModal = () => {
  openModal.value = false;
};
</script>
