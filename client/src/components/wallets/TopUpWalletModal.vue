<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleTopUpWallet">
      <div class="flex flex-col gap-4">
        <h1 class="text-xl font-bold capitalize dark:text-white">
          Top Up Wallet
        </h1>

        <div
          v-for="(currency, index) in props.walletCurrencies"
          :key="index"
          class="flex gap-4"
        >
          <text-input
            class="w-3/4 text-center"
            v-model.number="topUpAmounts[index]"
            placeholder="Amount"
          />
          <text-input
            class="w-1/4 rounded-md bg-slate-300 px-4 py-2 text-center text-sm font-bold text-gray-700 shadow-sm"
            v-model.trim="props.walletCurrencies[index]"
            disabled
          >
            {{ currency }}
          </text-input>
        </div>

        <div class="mt-4 flex justify-end gap-4">
          <action-button
            @click="closeModal"
            class="mb-4 inline-block rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300"
            text="Close"
          />
          <action-button
            class="mb-4 inline-block rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
            text="Submit"
          />
        </div>
      </div>
    </form>
  </modal>
</template>

<script lang="ts" setup>
import Modal from "@/components/Modal.vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { onMounted, ref, watch } from "vue";
import type { CurrencyAmount, TopUpWalletRequest } from "@/types/wallet";
import { TOP_UP_WALLET } from "@/api/wallet";
import type { GetUserBalanceCurrenciesResponse } from "@/types/balances";
import { GET_USER_BALANCE_CURRENCIES } from "@/api/balances";

const props = defineProps<{
  openModalTopUp: boolean;
  walletId: number | null;
  walletCurrencies: string[];
}>();

const isModalOpen = ref<boolean>(false);
const currencyAmountInputs = ref<CurrencyAmount[]>([
  { amount: null, currency: "" },
]);
const userBalanceCurrencies = ref<GetUserBalanceCurrenciesResponse[]>([]);
const topUpAmounts = ref<{ [index: number]: number }>({});

const emits = defineEmits(["closeModal", "formSubmitted"]);

onMounted(() => {
  getUserBalanceCurrencies();
});

watch(
  () => props.openModalTopUp,
  (newValue) => {
    if (props.openModalTopUp) {
      isModalOpen.value = newValue;
    }
  },
);

const handleTopUpWallet = async () => {
  try {
    const finalCurrencyAmount: CurrencyAmount[] = []; // Initialize an empty array to store the final data

    // Loop through the topUpAmounts object and create the final array
    for (const [index, amount] of Object.entries(topUpAmounts.value)) {
      finalCurrencyAmount.push({
        amount,
        currency: props.walletCurrencies[index],
      });
    }

    const body: TopUpWalletRequest = {
      currency_amount: finalCurrencyAmount,
    };

    const { status } = await TOP_UP_WALLET(props.walletId!, body);
    if (status === 204) {
      alert("Top Up Successful!");
      closeModal();
      emits("formSubmitted");
    }
  } catch (error: unknown) {
    alert(error);
  }
};

const getUserBalanceCurrencies = async () => {
  try {
    const { data, status } = await GET_USER_BALANCE_CURRENCIES();

    if (status === 200) {
      userBalanceCurrencies.value = data;
    }
  } catch (error: unknown) {
    alert(error);
  }
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};
</script>
