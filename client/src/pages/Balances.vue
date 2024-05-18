<template>
  <h1 class="text-xl font-bold">Balances</h1>
  <table class="table-auto">
    <thead>
      <tr>
        <th>Amount</th>
        <th>Currency</th>
        <th>Created At</th>
        <th>Updated At</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="balance in balances" class="text-center">
        <td>{{ balance.balance }}</td>
        <td>{{ balance.currency }}</td>
        <td>{{ formatDate(balance.createdAt) }}</td>
        <td>{{ formatDate(balance.updatedAt) }}</td>
        <td class="flex items-center justify-center gap-3">
          <action-button text="View" /><action-button
            text="Deposit / Withdraw"
          />
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts" setup>
import { GET_BALANCES } from "@/api/balances";
import type { Balance } from "@/types/balances";
import { onMounted, ref } from "vue";
import ActionButton from "@/components/ActionButton.vue";
import { format } from "date-fns";

let balances = ref([] as Balance[]);

onMounted(async () => {
  try {
    const { data, status } = await GET_BALANCES();

    balances.value = data.balances;
  } catch (error: unknown) {}
});

const formatDate = (dateString: string): string => {
  return format(new Date(dateString), "PPpp");
};
</script>
