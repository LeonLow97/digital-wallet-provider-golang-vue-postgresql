<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1 class="text-xl font-bold capitalize dark:text-white">
          Currency Exchange
        </h1>

        <div class="flex gap-4">
          <text-input
            v-model.number="fromAmount!"
            placeholder="Amount"
            class="w-3/4"
          />
          <text-input
            v-model="props.currency"
            disabled
            class="col-span-1 w-1/4 bg-slate-300 text-center font-bold uppercase"
          />
        </div>

        <div class="flex gap-4">
          <text-input
            v-model.number="toAmount!"
            placeholder="Amount"
            disabled
            class="col-span-1 w-3/4 bg-blue-100 text-center font-bold uppercase"
          />
          <select
            class="w-1/4 rounded-md border border-gray-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50"
            v-model.trim="selectedToCurrency"
          >
            <option disabled value="">To Currency</option>
            <option
              v-for="toCurrency in allowableToCurrencies"
              :key="toCurrency"
              :value="toCurrency"
            >
              <span>{{ toCurrency }}</span>
            </option>
          </select>
        </div>

        <div class="flex justify-end gap-4">
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
import { ref, watch, computed, onMounted } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { CURRENCY_EXCHANGE } from "@/api/balances";
import type { CURRENCY_EXCHANGE_REQUEST } from "@/types/balances";
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
    const body: CURRENCY_EXCHANGE_REQUEST = {
      from_amount: fromAmount.value!,
      to_currency: selectedToCurrency.value,
    };

    const { status } = await CURRENCY_EXCHANGE(body);

    if (status === 204) {
      emits("formSubmitted");
      emits("closeModal", true);
      toastStore.SUCCESS_TOAST("Successfully exchanged currency!");
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
