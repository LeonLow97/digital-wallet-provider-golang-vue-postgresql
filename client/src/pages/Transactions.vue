<template>
  <h1 class="text-xl font-bold tracking-wider">Transactions</h1>
  <div>{{ transactions }}</div>
</template>

<script lang="ts" setup>
import { GET_TRANSACTIONS } from "@/api/transactions";
import type { Transaction } from "@/types/transactions";
import { onMounted, ref } from "vue";

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
    alert(error?.response.data.message);
  }
};
</script>
