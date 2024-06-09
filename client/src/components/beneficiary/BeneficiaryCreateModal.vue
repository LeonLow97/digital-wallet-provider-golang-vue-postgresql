<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1
          class="text-center text-xl font-bold capitalize tracking-wider dark:text-white"
        >
          Create Beneficiary
        </h1>

        <div class="flex gap-4">
          <select
            class="rounded-md bg-white px-4 py-2 text-center text-sm text-gray-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50 dark:bg-gray-800 dark:text-gray-300"
            v-model.trim="mobileCountryCode"
          >
            <option value="+65">+ 65</option>
            <option value="+60">+ 60</option>
            <option value="+61">+ 61</option>
            <option value="+1">+ 1</option>
          </select>

          <text-input v-model.trim="mobileNumber" placeholder="Mobile Number" />
        </div>

        <em
          ><strong>NOTE:</strong> Please remove all whitespaces and dashes in
          mobile number.</em
        >

        <div class="flex justify-end gap-4">
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
import { ref, watch } from "vue";
import Modal from "@/components/Modal.vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { CREATE_BENEFICIARY } from "@/api/beneficiary";
import type { CREATE_BENEFICIARY_REQUEST } from "@/types/beneficiary";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const props = defineProps<{
  openModal: boolean;
}>();

watch(
  () => props.openModal,
  (newValue) => {
    if (props) {
      isModalOpen.value = newValue;
    }
  },
);

const emits = defineEmits(["closeModal", "formSubmitted"]);
const isModalOpen = ref<boolean>(false);
const mobileNumber = ref("");
const mobileCountryCode = ref("+65");

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};

const handleSubmit = async () => {
  try {
    const body: CREATE_BENEFICIARY_REQUEST = {
      mobile_country_code: mobileCountryCode.value,
      mobile_number: mobileNumber.value,
    };

    const { status } = await CREATE_BENEFICIARY(body);

    if (status === 201) {
      emits("formSubmitted");
      toastStore.SUCCESS_TOAST("Beneficiary Created Successfully!");
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  } finally {
    mobileNumber.value = "";
  }
};
</script>
