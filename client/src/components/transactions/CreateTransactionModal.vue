<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleTransfer">
      <div class="flex flex-col gap-4">
        <h1 class="text-xl font-bold capitalize dark:text-white">
          Make a Transfer
        </h1>

        <div class="mb-2">
          <label class="mb-2 block text-gray-700" for="wallet-type"
            >Select Currency for Transfer</label
          >
          <div class="flex gap-4">
            <text-input
              class="w-3/4"
              v-model.number="sourceAmount!"
              placeholder="Source Amount"
            />
            <select
              class="w-1/4 rounded-md border border-gray-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50"
              v-model.trim="selectedCurrency"
            >
              <option disabled value="">Select Currency for Transfer</option>
              <option
                v-for="currency in props.walletCurrencies"
                :key="currency"
                :value="currency"
              >
                {{ currency }}
              </option>
            </select>
          </div>
        </div>

        <div class="mb-4">
          <label class="mb-2 block text-gray-700" for="wallet-type"
            >Select Beneficiary</label
          >
          <select
            class="w-full rounded-md border border-gray-300 bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50"
            v-model.trim="selectedBeneficiary"
          >
            <option disabled value="">Select Beneficiary</option>
            <option
              v-for="beneficiary in beneficiaries"
              :key="beneficiary.beneficiaryMobileNumber"
              :value="{
                mobileCountryCode: beneficiary.beneficiaryMobileCountryCode,
                mobileNumber: beneficiary.beneficiaryMobileNumber,
              }"
            >
              <span
                >{{ beneficiary.beneficiaryFirstName }}
                {{ beneficiary.beneficiaryLastName }}, &nbsp;&nbsp;
                {{ beneficiary.beneficiaryMobileCountryCode }}
                {{ beneficiary.beneficiaryMobileNumber }}</span
              >
            </option>
          </select>
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
import { GET_BENEFICIARIES } from "@/api/beneficiary";
import type { GET_BENEFICIARY_RESPONSE } from "@/types/beneficiary";
import type { CreateTransactionRequest } from "@/types/transactions";
import { CREATE_TRANSACTION } from "@/api/transactions";
import { useToastStore } from "@/stores/toast";

const props = defineProps<{
  openTransactionModal: boolean;
  walletCurrencies: string[];
  walletId: number;
}>();
const toastStore = useToastStore();
const isModalOpen = ref<boolean>(false);
const beneficiaries = ref<GET_BENEFICIARY_RESPONSE[]>([]);
const selectedCurrency = ref();
const selectedBeneficiary = ref<{
  mobileCountryCode: string;
  mobileNumber: string | null;
}>({ mobileCountryCode: "", mobileNumber: null });
const sourceAmount = ref<number | null>(null);

onMounted(() => {
  getBeneficiaries();
});

watch(
  () => props.openTransactionModal,
  (newValue) => {
    if (props.openTransactionModal) {
      isModalOpen.value = newValue;
    }
  },
);

const emits = defineEmits(["closeTransactionModal", "formSubmitted"]);

const handleTransfer = async () => {
  try {
    const body: CreateTransactionRequest = {
      sender_wallet_id: props.walletId,
      source_currency: selectedCurrency.value,
      source_amount: sourceAmount.value!,
      beneficiary_mobile_country_code:
        selectedBeneficiary.value.mobileCountryCode,
      beneficiary_mobile_number: selectedBeneficiary.value.mobileNumber!,
    };

    const { status } = await CREATE_TRANSACTION(body);

    if (status === 201) {
      emits("formSubmitted", "Transaction Created");

      closeModal();
      clearData();
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};

const clearData = () => {
  selectedCurrency.value = "";
  selectedBeneficiary.value = { mobileCountryCode: "", mobileNumber: null };
  sourceAmount.value = null;
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeTransactionModal", true);
};

const getBeneficiaries = async () => {
  try {
    const { data, status } = await GET_BENEFICIARIES();

    beneficiaries.value = data.beneficiaries;
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
