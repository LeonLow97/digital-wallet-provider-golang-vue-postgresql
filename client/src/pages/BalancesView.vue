<template>
  <h1>Balance ID: {{ balance?.id }}</h1>
  <p>Balance: {{ balance?.balance }}</p>
  <p>Currency: {{ balance?.currency }}</p>
  <p>Created At: {{ balance?.createdAt }}</p>
  {{ balanceHistory }}
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GET_BALANCE, GET_BALANCE_HISTORY } from "@/api/balances";
import { useRoute, useRouter } from "vue-router";
import type {
  GetBalanceResponse,
  GetBalanceHistoryResponse,
} from "@/types/balances";

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
</script>
