<template>
  <div class="mx-auto mt-4 w-full">
    <!-- Balance Card -->
    <div class="flex justify-between">
      <div
        class="w-1/3 rounded-lg border border-gray-200 bg-white p-6 shadow-md"
      >
        <h2 class="mb-4 text-2xl font-bold text-gray-800">Balance</h2>
        <div class="flex flex-col gap-4">
          <div class="flex justify-between">
            <p class="text-lg text-gray-600">Balance:</p>
            <p class="text-lg text-blue-600">{{ balance?.balance }}</p>
          </div>
          <div class="flex justify-between">
            <p class="text-lg text-gray-600">Currency:</p>
            <p class="text-lg text-blue-600">{{ balance?.currency }}</p>
          </div>
          <div class="flex justify-between">
            <p class="text-lg text-gray-600">Created At:</p>
            <p class="text-lg italic text-gray-600">
              {{ formatDate(balance?.createdAt) }}
            </p>
          </div>
        </div>
      </div>

      <div>
        <router-link
          :to="{ name: 'Balances' }"
          class="text-lg text-blue-600 underline underline-offset-8 hover:text-blue-300"
          >&larr; Back to Balances</router-link
        >
      </div>
    </div>

    <!-- Balance History Table -->
    <h2 class="mb-2 mt-8 text-center text-xl font-bold">Balance History</h2>
    <div class="shadow-md">
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
            <td class="px-4 py-2">{{ history.amount }}</td>
            <td class="px-4 py-2">{{ history.currency }}</td>
            <td
              class="px-4 py-2 font-bold uppercase"
              :class="
                history.type.trim().toLowerCase() === 'deposit'
                  ? 'text-green-500'
                  : 'text-red-500'
              "
            >
              {{ history.type }}
            </td>
            <td class="px-4 py-2">{{ formatDate(history.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GET_BALANCE, GET_BALANCE_HISTORY } from "@/api/balances";
import { useRoute, useRouter } from "vue-router";
import type {
  GetBalanceResponse,
  GetBalanceHistoryResponse,
} from "@/types/balances";
import { format } from "date-fns";

const route = useRoute();
const router = useRouter();
let balance = ref<GetBalanceResponse | null>(null);
let balanceHistory = ref<GetBalanceHistoryResponse | null>(null);

onMounted(async () => {
  // ensure that params is a number
  const id = Number(route.params.id);
  if (isNaN(id)) {
    router.push({ name: "Balances" });
    return;
  }

  try {
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
  } catch (error: unknown) {
    alert(error);
  }
});

const formatDate = (dateString: string | undefined): string => {
  if (dateString) {
    return format(new Date(dateString), "PPpp");
  } else {
    return "Invalid Date Format";
  }
};
</script>
