<template>
  <h1 class="text-xl font-bold">Balances</h1>
  <table class="table-auto">
    <thead>
      <tr>
        <th>Amount</th>
        <th>Currency</th>
        <th>Created At</th>
        <th>Updated At</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="balance in balances">
        <td>{{ balance.balance }}</td>
        <td>{{ balance.currency }}</td>
        <td>{{ balance.createdAt }}</td>
        <td>{{ balance.updatedAt }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts" setup>
import { GET_BALANCES } from "@/api/balances";
import type { Balance } from "@/types/balances";
import { onMounted, ref } from "vue";

let balances = ref([] as Balance[]);

onMounted(async () => {
  try {
    const { data, status } = await GET_BALANCES();

    balances.value = data.balances;
  } catch (error: unknown) {}
});
</script>
