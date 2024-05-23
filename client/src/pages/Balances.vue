<template>
  <h1 class="text-xl font-bold">Balances</h1>
  <balances-table @form-submitted="formSubmitted" :balances="balances" />
</template>

<script lang="ts" setup>
import { GET_BALANCES } from "@/api/balances";
import type { Balance } from "@/types/balances";
import { onMounted, ref } from "vue";

import BalancesTable from "@/components/balances/BalancesTable.vue";

let balances = ref([] as Balance[]);

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
  } catch (error: unknown) {
    alert(error);
  }
};
</script>
