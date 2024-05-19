<template>
  <h1 class="text-xl font-bold">Balances</h1>
  <balances-table :balances="balances" />
</template>

<script lang="ts" setup>
import { GET_BALANCES } from "@/api/balances";
import type { Balance } from "@/types/balances";
import { onMounted, ref } from "vue";

import BalancesTable from "@/components/balances/BalancesTable.vue";

let balances = ref([] as Balance[]);

onMounted(async () => {
  try {
    const { data, status } = await GET_BALANCES();

    balances.value = data.balances;
  } catch (error: unknown) {
    alert(error);
  }
});
</script>
