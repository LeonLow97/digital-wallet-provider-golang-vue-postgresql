<template>
  <h1 class="mb-6 text-2xl font-bold">Enter Your MFA</h1>
  <text-input
    v-model.trim="mfaCode"
    maxlength="6"
    class="text-center"
    placeholder="6-digit code"
  />
  <action-button
    text="Submit"
    class="mt-4 w-full rounded-lg bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
    @click="handleVerifyMFA(email)"
  />
</template>

<script lang="ts" setup>
import { ref } from "vue";
import TextInput from "../TextInput.vue";
import ActionButton from "../ActionButton.vue";
import type { VerifyMfaRequest } from "@/types/user";
import { VERIFY_MFA } from "@/api/user";
import { useToastStore } from "@/stores/toast";

const toastStore = useToastStore();
const mfaCode = ref("");

defineProps({
  email: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["mfaVerified"]);

const handleVerifyMFA = async (email: string) => {
  try {
    const body: VerifyMfaRequest = {
      email,
      mfa_code: mfaCode.value,
    };

    const { status } = await VERIFY_MFA(body);

    if (status === 204) {
      emit("mfaVerified", true);
    }
  } catch (error: any) {
    toastStore.ERROR_TOAST(error?.response.data.message);
  }
};
</script>
