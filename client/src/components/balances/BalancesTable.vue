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
            @click="handleDeposit(balance)"
          />
          <action-button
            text="Withdraw"
            class="rounded-lg border-2 border-pink-500 px-2 text-center text-pink-500 transition hover:border-pink-200 hover:text-pink-300"
            @click="handleWithdraw(balance)"
          />
        </td>
        <td class="px-4 py-2">{{ formatDate(balance.createdAt) }}</td>
        <td class="px-4 py-2">{{ formatDate(balance.updatedAt) }}</td>
      </tr>
    </tbody>
  </table>

  <balance-modal
    @close-modal="closeModal"
    @form-submitted="formSubmitted"
    class="text-left"
    :open-modal="openModal"
    :action-type="actionType"
    :balance="selectedBalance!"
  />
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { format } from "date-fns";
import { useRouter } from "vue-router";

import type { Balance } from "@/types/balances";
import ActionButton from "../ActionButton.vue";
import BalanceModal from "@/components/balances/BalanceModal.vue";

defineProps<{ balances: Balance[] }>();

const emits = defineEmits(["formSubmitted"]);
const router = useRouter();
let actionType = ref("");
const openModal = ref(false);
const selectedBalance = ref<Balance>({
  id: 0,
  balance: 0,
  currency: "",
  createdAt: "",
  updatedAt: "",
});

const formatDate = (dateString: string | undefined): string | undefined => {
  if (dateString) {
    return format(new Date(dateString), "PPpp");
  }
};

const handleClickViewBalances = (balanceId: number) => {
  router.push({ name: "Balance", params: { id: balanceId } });
};

const handleDeposit = (balance: Balance) => {
  actionType.value = "deposit";
  openModal.value = true;
  selectedBalance.value = balance;
};

const handleWithdraw = (balance: Balance) => {
  actionType.value = "withdraw";
  openModal.value = true;
  selectedBalance.value = balance;
};

const closeModal = () => {
  openModal.value = false;
};

const formSubmitted = () => {
  emits("formSubmitted");
};
</script>
