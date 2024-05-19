<template>
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
          <action-button
            text="View"
            class="rounded-lg border-2 border-blue-500 px-2 text-center text-blue-500 transition hover:border-blue-200 hover:text-blue-300"
            @click="handleClickViewBalances(balance.id)"
          />
          <action-button
            text="Deposit / Withdraw"
            class="rounded-lg border-2 border-blue-500 px-2 text-center text-blue-500 transition hover:border-blue-200 hover:text-blue-300"
          />
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts" setup>
import type { Balance } from "@/types/balances";
import { format } from "date-fns";
import ActionButton from "../ActionButton.vue";
import { useRouter } from "vue-router";

defineProps<{ balances: Balance[] }>();

const router = useRouter();

const formatDate = (dateString: string): string => {
  return format(new Date(dateString), "PPpp");
};

const handleClickViewBalances = (balanceId: number) => {
  router.push({ name: "BalancesView", params: { id: balanceId } });
};
</script>
