<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1 class="text-xl font-bold dark:text-white">
          Are you sure you want to
          <span class="uppercase text-blue-500">{{
            props.actionDelete ? "delete" : "reactivate"
          }}</span>
          this beneficiary?
        </h1>
        <div class="flex justify-end gap-4">
          <action-button
            @click="closeModal"
            class="mb-4 inline-block rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300"
            text="No"
          />
          <action-button
            class="mb-4 inline-block rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
            text="Yes"
          />
        </div>
      </div>
    </form>
  </modal>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import Modal from "@/components/Modal.vue";
import ActionButton from "@/components/ActionButton.vue";
import { DELETE_BENEFICIARY } from "@/api/beneficiary";
import type { DELETE_BENEFICIARY_REQUEST } from "@/types/beneficiary";

const props = defineProps<{
  openModal: boolean;
  beneficiaryId: number;
  actionDelete: boolean;
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

const closeModal = () => {
  isModalOpen.value = false;
  emits("closeModal", true);
};

const handleSubmit = async () => {
  try {
    const body: DELETE_BENEFICIARY_REQUEST = {
      is_deleted: props.actionDelete ? 1 : 0,
      beneficiary_id: props.beneficiaryId,
    };

    const { status } = await DELETE_BENEFICIARY(body);

    if (status === 204) {
      alert("Updated Successfully!");
      emits("formSubmitted");
    }
  } catch (error: any) {
    alert(error?.response.data.message);
  }
};
</script>
