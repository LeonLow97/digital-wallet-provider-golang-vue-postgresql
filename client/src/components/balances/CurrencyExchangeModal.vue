<template>
  <modal @close-overlay="closeModal" modal-width="1/2" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col">
        <h1
          class="mb-4 text-center text-xl font-bold capitalize tracking-wider dark:text-white"
        >
          Currency Exchange
        </h1>

        <div class="mb-2 font-bold text-gray-700">Please Select One:</div>

        <div class="mb-4 flex">
          <div
            @click="handleClickedAmountToSend"
            class="cursor-pointer px-4 py-4 uppercase hover:bg-gray-100 tracking-wide"
            :class="
              showAmountToSend ? 'border-b-2 border-blue-500 text-blue-500' : ''
            "
          >
            <span>Amount to Send</span>
          </div>
          <div
            class="cursor-pointer px-4 py-4 uppercase hover:bg-gray-100 tracking-wide"
            @click="handleClickedAmountToReceive"
            :class="
              showAmountToReceive
                ? 'border-b-2 border-blue-500 text-blue-500'
                : ''
            "
          >
            <span>Amount to Receive</span>
          </div>
        </div>

        <div
          class="grid grid-cols-5 items-center gap-4"
          v-if="showAmountToSend"
        >
          From:
          <text-input
            v-model.number="fromAmount!"
            placeholder="Amount"
            class="col-span-3"
          />
          <text-input
            v-model="props.currency"
            disabled
            class="col-span-1 bg-slate-200 text-center font-bold uppercase text-slate-500"
          />
        </div>

        <div class="mt-4 grid grid-cols-5 items-center gap-4">
          To:
          <text-input
            v-if="showAmountToReceive"
            v-model.number="toAmount!"
            placeholder="Amount"
            class="col-span-3"
          />
          <select
            class="cols-span-1 text-md rounded-md border border-gray-300 px-4 py-3 text-center shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50"
            v-model.trim="selectedToCurrency"
          >
            <option disabled value="">Select Currency</option>
            <option
              v-for="toCurrency in allowableToCurrencies"
              :key="toCurrency"
              :value="toCurrency"
            >
              <span>{{ toCurrency }}</span>
            </option>
          </select>
        </div>

        <div class="mt-8 flex justify-center">
          <action-button
            type="button"
            @click="handleExchangeCalculation"
            class="mb-8 inline-block rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
            text="Calculate Exchange"
          />
        </div>

        <div class="bg-gray-100 p-4" v-if="exchangeCalculated">
          <div class="mb-2 font-bold uppercase tracking-wider">
            Indicative Exchange:
          </div>
          <div class="mb-1 flex gap-6">
            <p>From:</p>
            <div class="flex gap-3">
              <span>{{ previewExchangeResponse?.fromAmount?.toFixed(2) }}</span>
              <span>{{ previewExchangeResponse?.fromCurrency }}</span>
            </div>
          </div>
          <div class="flex gap-6">
            <p>To:</p>
            <div class="flex gap-3">
              <span>{{ previewExchangeResponse?.toAmount?.toFixed(2) }}</span>
              <span>{{ previewExchangeResponse?.toCurrency }}</span>
            </div>
          </div>
        </div>

        <div class="mt-8 flex justify-end gap-4">
          <action-button
            @click="closeModal"
            class="mb-4 inline-block rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300"
            text="Close"
          />
          <action-button
            v-if="exchangeCalculated"
            class="500 mb-4 inline-block rounded-lg px-4 py-2 text-center text-white transition"
            :class="
              !exchangeCalculated
                ? 'bg-blue-300'
                : 'bg-blue-500 hover:bg-blue-400'
            "
            text="Submit"
            :disabled="!exchangeCalculated"
          />
        </div>
      </div>
    </form>
  </modal>
</template>

<script lang="ts" setup>
import Modal from "@/components/Modal.vue";
import { ref, watch, onMounted } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { CURRENCY_EXCHANGE, PREVIEW_EXCHANGE } from "@/api/balances";
import type {
  CURRENCY_EXCHANGE_REQUEST,
  PreviewExchangeRequest,
  PreviewExchangeResponse,
} from "@/types/balances";
import { useToastStore } from "@/stores/toast";

const props = defineProps<{
  openModal: boolean;
  currency: string;
}>();
const toastStore = useToastStore();
const isModalOpen = ref<boolean>(false);
const fromAmount = ref<number | null>(null);
const toAmount = ref<number | null>(null); // use computed property instead
const selectedToCurrency = ref("");
const allowableToCurrencies = ref<string[]>([]);

const showAmountToSend = ref(true);
const showAmountToReceive = ref(false);

const previewExchangeResponse = ref<PreviewExchangeResponse | null>(null);

// state to keep track whether user has clicked "Calculate Exchange" button
const exchangeCalculated = ref(false);

const emits = defineEmits(["closeModal", "formSubmitted"]);

onMounted(() => {
  // TODO: retrieve this from backend instead
  let toCurrencies = ["USD", "SGD", "AUD", "MYR"];
  allowableToCurrencies.value = toCurrencies.filter(
    (item) => item !== props.currency,
  );
});

watch(
  () => props.openModal,
  (newValue) => {
    if (props) {
      isModalOpen.value = newValue;
    }
  },
);

const handleSubmit = async () => {
  try {
    if (!exchangeCalculated.value) {
      return;
    }

    const body: CURRENCY_EXCHANGE_REQUEST = {
      from_amount: fromAmount.value!,
      to_currency: selectedToCurrency.value,
    };

    const { status } = await CURRENCY_EXCHANGE(body);

    if (status === 204) {
      emits("formSubmitted");
      emits("closeModal", true);

      clearData();
      toastStore.SUCCESS_TOAST("Successfully exchanged currency!");
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const handleClickedAmountToSend = () => {
  showAmountToSend.value = true;
  showAmountToReceive.value = !showAmountToSend.value;
};

const handleClickedAmountToReceive = () => {
  showAmountToReceive.value = true;
  showAmountToSend.value = !showAmountToReceive.value;
};

const handleExchangeCalculation = async () => {
  exchangeCalculated.value = true;

  try {
    let body: PreviewExchangeRequest = {
      action_type: null,
      from_amount: 0,
      from_currency: "",
      to_amount: 0,
      to_currency: "",
    };

    if (showAmountToSend.value) {
      body.action_type = "amountToSend";
      body.from_amount = fromAmount.value!;
      body.from_currency = props.currency;
      body.to_currency = selectedToCurrency.value;
    }
    if (showAmountToReceive.value) {
      body.action_type = "amountToReceive";
      body.from_currency = props.currency;
      body.to_amount = toAmount.value!;
      body.to_currency = selectedToCurrency.value;
    }

    const { data, status } = await PREVIEW_EXCHANGE(body);

    if (status === 200) {
      previewExchangeResponse.value = data;

      if (showAmountToReceive.value) {
        fromAmount.value = data.fromAmount;
      }
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const clearData = () => {
  fromAmount.value = null;
  selectedToCurrency.value = "";
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};
</script>
