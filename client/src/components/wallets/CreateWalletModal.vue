<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleCreateWallet">
      <div class="flex flex-col gap-4">
        <h1
          class="text-center text-xl font-bold capitalize tracking-wider dark:text-white"
        >
          Create Wallet
        </h1>

        <div class="mb-4">
          <label
            class="mb-2 block text-gray-700 dark:text-gray-300"
            for="wallet-type"
            >Wallet Type:</label
          >
          <select
            class="w-full rounded-md border border-slate-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50 dark:border-none dark:bg-gray-800 dark:text-gray-300"
            v-model.trim="selectedWalletTypeId"
          >
            <option disabled value="0">Select Wallet Type</option>
            <option
              v-for="walletType in walletTypes"
              :key="walletType.id"
              :value="walletType.id"
            >
              <span>{{ capitalizeFirstLetter(walletType.walletType) }}</span>
            </option>
          </select>
        </div>

        <div
          class="flex gap-4"
          v-for="(item, index) in currencyAmountInputs"
          :key="index"
        >
          <text-input v-model.number="item.amount!" placeholder="Amount" />
          <select
            class="w-1/3 rounded-md border border-slate-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50 dark:border-none dark:bg-gray-800 dark:text-gray-300"
            v-model.trim="item.currency"
          >
            <option disabled value="">Select Currency</option>
            <option
              v-for="item in userBalanceCurrencies"
              :key="item.currency"
              :value="item.currency"
            >
              {{ item.currency }}
            </option>
          </select>
        </div>
        <button
          @click.prevent="addCurrency"
          class="text-blue-800 underline underline-offset-8 transition hover:text-blue-400 dark:text-blue-300 dark:hover:text-blue-500"
          v-if="currencyAmountInputs.length !== userBalanceCurrencies.length"
        >
          Add more currencies
        </button>

        <div class="mt-4 flex justify-end gap-4">
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
      </div>
    </form>
  </modal>
</template>

<script lang="ts" setup>
import Modal from "@/components/Modal.vue";
import { onMounted, ref, watch } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import type {
  CreateWalletRequest,
  CurrencyAmount,
  GetWalletTypesResponse,
} from "@/types/wallet";
import { CREATE_WALLET, GET_WALLET_TYPES } from "@/api/wallet";
import { GET_USER_BALANCE_CURRENCIES } from "@/api/balances";
import type { GetUserBalanceCurrenciesResponse } from "@/types/balances";
import { useToastStore } from "@/stores/toast";

const props = defineProps<{
  openModal: boolean;
}>();
const toastStore = useToastStore();
const selectedWalletTypeId = ref(0);
const currencyAmountInputs = ref<CurrencyAmount[]>([
  { amount: null, currency: "" },
]);

const walletTypes = ref<GetWalletTypesResponse[]>([]);
const userBalanceCurrencies = ref<GetUserBalanceCurrenciesResponse[]>([]);
const isModalOpen = ref<boolean>(false);

const emits = defineEmits(["closeModal", "formSubmitted"]);

onMounted(() => {
  // TODO: Make this run asynchronously
  getWalletTypes();
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

const addCurrency = () => {
  if (currencyAmountInputs.value.length < userBalanceCurrencies.value.length) {
    currencyAmountInputs.value.push({ amount: null, currency: "" });
  }
};

const capitalizeFirstLetter = (word: string) => {
  return word.charAt(0).toUpperCase() + word.slice(1);
};

const getWalletTypes = async () => {
  try {
    const { data, status } = await GET_WALLET_TYPES();

    if (status === 200) {
      walletTypes.value = data;
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
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

const handleCreateWallet = async () => {
  try {
    // if currency or amount not provided, remove that currencyAmount
    const filteredCurrencyAmount = currencyAmountInputs.value.filter(
      (currencyAmount) =>
        currencyAmount.currency !== "" || currencyAmount.amount !== null,
    );

    const body: CreateWalletRequest = {
      wallet_type_id: selectedWalletTypeId.value,
      currency_amount: filteredCurrencyAmount,
    };

    const { status } = await CREATE_WALLET(body);

    if (status === 201) {
      toastStore.SUCCESS_TOAST("Created Wallet Successfully!");

      selectedWalletTypeId.value = 0;
      emits("formSubmitted");
      closeModal();
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
  currencyAmountInputs.value = [{ amount: null, currency: "" }];
};
</script>
