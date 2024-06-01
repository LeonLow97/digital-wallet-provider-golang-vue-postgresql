<template>
  <div>
    <h1 class="text-xl font-bold tracking-wider mb-4">Transactions</h1>
    <div class="overflow-x-auto">
      <table class="bg-white border border-gray-200">
        <thead>
          <tr class="bg-gray-200 text-sm leading-normal">
            <th class="py-3 px-6 text-center">Sender Username</th>
            <th class="py-3 px-6 text-center">Sender Mobile</th>
            <th class="py-3 px-6 text-center">Beneficiary Username</th>
            <th class="py-3 px-6 text-center">Beneficiary Mobile</th>
            <th class="py-3 px-6 text-center">Source Amount</th>
            <th class="py-3 px-6 text-center">Source Currency</th>
            <th class="py-3 px-6 text-center">Destination Amount</th>
            <th class="py-3 px-6 text-center">Destination Currency</th>
            <th class="py-3 px-6 text-center">Source of Transfer</th>
            <th class="py-3 px-6 text-center">Status</th>
            <th class="py-3 px-6 text-center">Created At</th>
          </tr>
        </thead>
        <tbody class=" text-sm">
          <tr v-for="transaction in transactions" :key="transaction.created_at" class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-6 text-center">{{ transaction.sender_username }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.sender_mobile_number }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.beneficiary_username }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.beneficiary_mobile_number }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.source_amount }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.source_currency }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.destination_amount }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.destination_currency }}</td>
            <td class="py-3 px-6 text-center">{{ transaction.source_of_transfer }}</td>
            <td class="py-3 px-6 text-center" :class="transaction.status === 'COMPLETED' ? 'text-green-500' : 'text-red-500'">{{ transaction.status }}</td>
            <td class="py-3 px-6 text-center">{{ new Date(transaction.created_at).toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { GET_TRANSACTIONS } from "@/api/transactions";
import type { Transaction } from "@/types/transactions";
import { onMounted, ref } from "vue";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const transactions = ref<Transaction[]>([]);

onMounted(() => {
  getTransactions();
});

const getTransactions = async () => {
  try {
    const { data, status } = await GET_TRANSACTIONS();

    if (status === 200) {
      transactions.value = data;
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
