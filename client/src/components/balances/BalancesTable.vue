<template>
  <table class="w-full table-fixed text-sm shadow-md">
    <thead class="bg-gray-100">
      <tr class="text-lg font-bold">
        <th class="px-4 py-2">Amount</th>
        <th class="px-4 py-2">Currency</th>
        <th class="px-4 py-2">Actions</th>
        <th class="px-4 py-2">Created At</th>
        <th class="px-4 py-2">Updated At</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="balance in balances" class="text-center hover:bg-gray-50">
        <td class="px-4 py-2">{{ balance.balance }}</td>
        <td class="px-4 py-2">{{ balance.currency }}</td>
        <td class="flex items-center justify-center gap-3 px-4 py-2">
          <action-button
            text="View"
            class="rounded-lg border-2 border-blue-500 px-2 text-center text-blue-500 transition hover:border-blue-200 hover:text-blue-300"
            @click="handleClickViewBalances(balance.id)"
          />
          <action-button
            text="Deposit"
            class="rounded-lg border-2 border-green-500 px-2 text-center text-green-500 transition hover:border-green-200 hover:text-green-300"
            @click="handleDeposit"
          />
          <action-button
            text="Withdraw"
            class="rounded-lg border-2 border-orange-500 px-2 text-center text-orange-500 transition hover:border-orange-200 hover:text-orange-300"
            @click="handleWithdraw"
          />
        </td>
        <td class="px-4 py-2">{{ formatDate(balance.createdAt) }}</td>
        <td class="px-4 py-2">{{ formatDate(balance.updatedAt) }}</td>
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
  router.push({ name: "Balance", params: { id: balanceId } });
};

const handleDeposit = () => {
  alert("deposit");
};

const handleWithdraw = () => {
  alert("withdraw");
};
</script>
