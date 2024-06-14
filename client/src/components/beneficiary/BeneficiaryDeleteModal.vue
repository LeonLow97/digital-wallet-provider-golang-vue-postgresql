<template>
  <modal @close-overlay="closeModal" modal-width="1/3" v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <h1 class="text-xl font-bold tracking-wide dark:text-white">
          Are you sure you want to
          <span class="uppercase text-blue-500 dark:text-blue-400">{{
            props.actionDelete ? "delete" : "reactivate"
          }}</span>
          this beneficiary?
        </h1>
        <div class="flex justify-end gap-4">
          <action-button
            @click="closeModal"
            class="mb-4 inline-block rounded-lg border border-blue-500 px-4 py-2 text-center text-blue-500 transition hover:border-blue-300 hover:text-blue-300 dark:border-blue-300 dark:text-blue-300 dark:hover:border-blue-500 dark:hover:text-blue-500"
            text="No"
          />
          <action-button
            class="mb-4 inline-block rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
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
import type { DeleteBeneficiaryRequest } from "@/types/beneficiary";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
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
    const body: DeleteBeneficiaryRequest = {
      is_deleted: props.actionDelete ? 1 : 0,
      beneficiary_id: props.beneficiaryId,
    };

    const { status } = await DELETE_BENEFICIARY(body);

    if (status === 204) {
      toastStore.SUCCESS_TOAST(
        `${props.actionDelete ? "Deleted" : "Reactivated"} Successfully!`,
      );
      emits("formSubmitted");
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
