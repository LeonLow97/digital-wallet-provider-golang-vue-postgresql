<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleTopUpWallet">
      <h1
        class="mb-6 text-center text-xl font-bold capitalize tracking-wider dark:text-white"
      >
        {{ props.actionType }} Wallet
      </h1>

      <div
        v-for="(currency, index) in props.walletCurrencies"
        :key="index"
        class="mb-4 grid grid-cols-3 gap-4"
      >
        <text-input
          class="col-span-2 text-center"
          v-model.number="topUpAmounts[index]"
          placeholder="Amount"
        />
        <text-input
          class="col-span-1 rounded-md bg-slate-200 px-4 py-2 text-center text-sm font-bold text-slate-500 shadow-sm"
          v-model.trim="props.walletCurrencies[index]"
          disabled
        >
          {{ currency }}
        </text-input>
      </div>

      <div class="mt-8 flex justify-end gap-4">
        <action-button
          @click="closeModal"
          class="mb-4 inline-block rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300 dark:border-blue-300 dark:text-blue-300 dark:hover:border-blue-500 dark:hover:text-blue-500"
          text="Close"
        />
        <action-button
          class="mb-4 inline-block rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
          text="Submit"
        />
      </div>
    </form>
  </modal>
</template>

<script lang="ts" setup>
import Modal from "@/components/Modal.vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { onMounted, ref, watch } from "vue";
import type { CurrencyAmount, WalletExchangesRequest } from "@/types/wallet";
import { TOP_UP_WALLET, CASH_OUT_WALLET } from "@/api/wallet";
import type { GetUserBalanceCurrenciesResponse } from "@/types/balances";
import { GET_USER_BALANCE_CURRENCIES } from "@/api/balances";
import { useToastStore } from "@/stores/toast";

const props = defineProps<{
  openModal: boolean;
  walletId: number | null;
  walletCurrencies: string[];
  actionType: string;
}>();
const toastStore = useToastStore();
const isModalOpen = ref<boolean>(false);

const userBalanceCurrencies = ref<GetUserBalanceCurrenciesResponse[]>([]);
const topUpAmounts = ref<{ [index: number]: number }>({});

const emits = defineEmits(["closeModal", "formSubmitted"]);

onMounted(() => {
  getUserBalanceCurrencies();
});

watch(
  () => props.openModal,
  (newValue) => {
    if (props.openModal) {
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

    const body: WalletExchangesRequest = {
      currency_amount: finalCurrencyAmount,
    };

    let responseStatus: number = 0;

    if (props.actionType === "Top Up") {
      const { status } = await TOP_UP_WALLET(props.walletId!, body);
      responseStatus = status;
    } else if (props.actionType === "Cash Out") {
      const { status } = await CASH_OUT_WALLET(props.walletId!, body);
      responseStatus = status;
    }

    if (responseStatus === 204) {
      closeModal();
      emits("formSubmitted", props.actionType);
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  } finally {
    // clean data
    topUpAmounts.value = {};
  }
};

const getUserBalanceCurrencies = async () => {
  try {
    const { data, status } = await GET_USER_BALANCE_CURRENCIES();

    if (status === 200) {
      userBalanceCurrencies.value = data;
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};
</script>
