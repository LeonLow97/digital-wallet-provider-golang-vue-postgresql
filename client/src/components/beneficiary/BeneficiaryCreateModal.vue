<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1 class="text-xl font-bold capitalize dark:text-white">
          Create Beneficiary
        </h1>
        <text-input v-model.trim="mobileNumber" placeholder="Mobile Number" />

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
import { ref, watch } from "vue";
import Modal from "@/components/Modal.vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import { CREATE_BENEFICIARY } from "@/api/beneficiary";
import type { CREATE_BENEFICIARY_REQUEST } from "@/types/beneficiary";

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

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};

const handleSubmit = async () => {
  try {
    const body: CREATE_BENEFICIARY_REQUEST = {
      mobile_number: mobileNumber.value,
    };

    const { status } = await CREATE_BENEFICIARY(body);

    if (status === 201) {
      emits("formSubmitted");
      alert("Beneficiary created successfully!");
    }
  } catch (error: unknown) {
    alert(error);
  } finally {
    mobileNumber.value = "";
  }
};
</script>
