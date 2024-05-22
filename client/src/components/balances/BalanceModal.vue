<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1
          :class="
            props.actionType.trim().toLowerCase() === 'deposit'
              ? 'text-green-500'
              : 'text-red-500'
          "
          class="text-xl font-bold capitalize dark:text-white"
        >
          {{ props.actionType }}
        </h1>
        <div class="grid grid-cols-3 gap-4">
          <text-input
            v-model.number="amount"
            placeholder="Amount"
            class="col-span-2"
          />
          <text-input
            v-model="currency"
            placeholder="Currency"
            disabled
            class="col-span-1 bg-slate-300 text-center font-bold uppercase"
          />
        </div>
        <p>Current Balance: {{ props.balance?.balance }}</p>
        <p v-if="amount">
          Final Balance:
          <span
            :class="
              props.actionType.trim().toLowerCase() === 'deposit'
                ? 'text-green-500'
                : 'text-red-500'
            "
            class="font-bold"
            >{{ finalBalance }}</span
          >
        </p>
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
import { computed, ref, watch } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { DEPOSIT, WITHDRAW } from "@/api/balances";
import type { DEPOSIT_REQUEST, WITHDRAW_REQUEST } from "@/types/balances";
import type { GENERIC_STATUS_RESPONSE } from "@/types/generic";

const props = defineProps<{
  actionType: string;
  openModal: boolean;
  balance: {
    id: number;
    balance: number;
    currency: string;
    createdAt: string;
  } | null;
}>();

const isModalOpen = ref<boolean>(false);
const amount = ref<number>(0);
const currency = ref("");

const emits = defineEmits(["closeModal", "formSubmitted"]);

watch(
  () => props.openModal,
  (newValue) => {
    if (props) {
      isModalOpen.value = newValue;
    }
  },
);

watch(
  () => props.balance?.currency,
  (newValue) => {
    if (newValue !== undefined) {
      currency.value = newValue;
    }
  },
);

const finalBalance = computed(() => {
  let finalValue = 0; // Initialize with default value

  if (props.balance?.balance && amount.value) {
    if (props.actionType?.toLowerCase() === "withdraw") {
      finalValue = props.balance?.balance - amount.value;
    } else if (props.actionType?.toLowerCase() === "deposit") {
      finalValue = props.balance?.balance + amount.value;
    }
  }

  return finalValue <= 0 ? 0 : finalValue;
});

const handleSubmit = async () => {
  try {
    const body: DEPOSIT_REQUEST = {
      amount: amount.value,
      currency: currency.value,
    };

    console.log(props.actionType)

    let response: GENERIC_STATUS_RESPONSE;
    if (props.actionType?.trim().toLowerCase() === "deposit") {
      response = await DEPOSIT(body);
    } else {
      response = await WITHDRAW(body);
    }

    const { status } = response;
    if (status === 204) {
      alert("Submitted Successfully!");
      closeModal();
      emits("formSubmitted");
    }
  } catch (error: unknown) {
    alert(error);
  }
};

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);

  amount.value = 0;
};
</script>
