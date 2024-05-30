<template>
  <h1 class="text-xl font-bold tracking-wider">Beneficiary</h1>
  <div class="flex items-center justify-between">
    <div class="p-4 text-center tracking-wider text-blue-600 shadow-md">
      You have
      <span
        class="font-bold uppercase"
        :class="showDeletedBeneficiaries ? 'text-red-500' : 'text-green-500'"
        >{{ filterBeneficiaries ? filterBeneficiaries.length : 0 }}
        {{ showDeletedBeneficiaries ? "inactive" : "active" }}</span
      >
      {{ filterBeneficiaries?.length === 1 ? "beneficiary" : "beneficiaries" }}
    </div>
    <div class="flex gap-8">
      <action-button
        class="hover rounded-lg bg-blue-500 px-4 py-2 text-white transition hover:bg-blue-300"
        text="Create Beneficiary"
        @click="handleCreateBeneficiary"
      />
      <label for="" class="flex items-center gap-2">
        <input
          type="checkbox"
          class="h-6 w-6 accent-blue-500"
          v-model="showDeletedBeneficiaries"
        />
        Show Deleted Beneficiaries
      </label>
    </div>
  </div>

  <h2 class="mt-4 text-lg font-bold">My Beneficiaries</h2>
  <div class="overflow-x-auto shadow-md">
    <table class="w-full text-center">
      <thead>
        <tr class="bg-gray-100">
          <th class="px-4 py-2 text-lg font-bold">First Name</th>
          <th class="px-4 py-2 text-lg font-bold">Last Name</th>
          <th class="px-4 py-2 font-bold">Email</th>
          <th class="px-4 py-2 font-bold">Username</th>
          <th class="px-4 py-2 font-bold">Mobile Number</th>
          <th class="px-4 py-2 font-bold">Status</th>
          <th class="px-4 py-2 font-bold">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="beneficiary in filterBeneficiaries"
          :key="beneficiary.beneficiaryID"
          class="bg-white hover:bg-gray-50"
        >
          <td class="px-4 py-2">
            {{ beneficiary.beneficiaryFirstName }}
          </td>
          <td class="px-4 py-2">
            {{ beneficiary.beneficiaryLastName }}
          </td>
          <td class="px-4 py-2">
            {{ beneficiary.beneficiaryEmail }}
          </td>
          <td class="px-4 py-2">
            {{ beneficiary.beneficiaryUsername }}
          </td>
          <td class="px-4 py-2">
            {{ beneficiary.beneficiaryMobileCountryCode }}
            {{ beneficiary.beneficiaryMobileNumber }}
          </td>
          <td class="px-4 py-2">
            <span v-if="beneficiary.active === 1" class="text-green-500"
              >Active</span
            >
            <span v-else class="text-red-500">Inactive</span>
          </td>
          <td>
            <button @click="handleDeleteBeneficiary(beneficiary)">
              <svg-icon
                type="mdi"
                :path="
                  showDeletedBeneficiaries
                    ? mdiAccountReactivateOutline
                    : mdiTrashCanOutline
                "
              />
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="!beneficiaries" class="p-4 text-center">
      No Beneficiaries Found!
    </div>
  </div>

  <beneficiary-create-modal
    @close-modal="closeCreateBeneficiaryModal"
    @form-submitted="createBeneficiaryFormSubmitted"
    :open-modal="openCreateBeneficiaryModal"
  />

  <beneficiary-delete-modal
    @close-modal="closeDeleteBeneficiaryModal"
    @form-submitted="deleteBeneficiaryFormSubmitted"
    :open-modal="openDeleteBeneficiaryModal"
    :beneficiary-id="deleteBeneficiaryId"
    :action-delete="!showDeletedBeneficiaries"
  />
</template>

<script lang="ts" setup>
import ActionButton from "@/components/ActionButton.vue";
import { GET_BENEFICIARIES } from "@/api/beneficiary";
import type { GET_BENEFICIARY_RESPONSE } from "@/types/beneficiary";
import { onMounted, ref, computed } from "vue";
import BeneficiaryCreateModal from "@/components/beneficiary/BeneficiaryCreateModal.vue";
import BeneficiaryDeleteModal from "@/components/beneficiary/BeneficiaryDeleteModal.vue";
import { mdiTrashCanOutline, mdiAccountReactivateOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";

const beneficiaries = ref<GET_BENEFICIARY_RESPONSE[]>([]);
const deleteBeneficiaryId = ref<number>(0);

const openCreateBeneficiaryModal = ref(false);
const openDeleteBeneficiaryModal = ref(false);
const showDeletedBeneficiaries = ref(false);

onMounted(async () => {
  getBeneficiaries();
});

const filterBeneficiaries = computed(() => {
  if (beneficiaries.value) {
    if (showDeletedBeneficiaries.value) {
      return beneficiaries?.value.filter(
        (beneficiary: GET_BENEFICIARY_RESPONSE) => beneficiary.isDeleted === 1,
      );
    }
    return beneficiaries?.value.filter(
      (beneficiary: GET_BENEFICIARY_RESPONSE) => beneficiary.isDeleted === 0,
    );
  }
});

const getBeneficiaries = async () => {
  try {
    const { data, status } = await GET_BENEFICIARIES();

    beneficiaries.value = data.beneficiaries;
  } catch (error: unknown) {
    alert(error);
  }
};

const handleCreateBeneficiary = async () => {
  openCreateBeneficiaryModal.value = true;
};

const handleDeleteBeneficiary = async (
  beneficiary: GET_BENEFICIARY_RESPONSE,
) => {
  openDeleteBeneficiaryModal.value = true;
  deleteBeneficiaryId.value = beneficiary.beneficiaryID;
};

const closeCreateBeneficiaryModal = () => {
  openCreateBeneficiaryModal.value = false;
};

const closeDeleteBeneficiaryModal = () => {
  openDeleteBeneficiaryModal.value = false;
};

const createBeneficiaryFormSubmitted = () => {
  openCreateBeneficiaryModal.value = false;
  getBeneficiaries();
};

const deleteBeneficiaryFormSubmitted = () => {
  openDeleteBeneficiaryModal.value = false;
  getBeneficiaries();
};
</script>
