<template>
  <div class="flex items-center justify-center gap-12">
    <div class="w-1/2">
      <h1 class="mb-6 text-2xl font-bold">Configure Your MFA</h1>

      <div class="mb-6">
        <p class="mb-2">
          To enhance the security of your account, we require you to set up
          Multi-Factor Authentication (MFA) using the Google Authenticator app.
        </p>
      </div>

      <div class="mb-6">
        <p class="mb-2 font-bold">Follow these steps to configure MFA:</p>
        <ol class="list-decimal pl-6">
          <li class="mb-2">
            Add the Google Chrome Extension
            <a
              href="https://chromewebstore.google.com/detail/authenticator/bhghoamapcdpbohphigoooaddinpkbai?pli=1"
              class="text-blue-600 underline dark:text-blue-400"
              target="_blank"
              rel="noopener"
              >Google Authenticator Extension</a
            >
            on your device.
          </li>
          <li class="mb-2">
            Open the Google Authenticator app and tap on the '+' icon to add a
            new account.
          </li>
          <li class="mb-2">
            Select the option to scan a QR code and point your device's camera
            at the QR code displayed on the right.
          </li>
          <li class="mb-2">
            Once scanned, the app will generate a unique code for your account.
            Enter this code to complete the setup.
          </li>
        </ol>
      </div>

      <div>
        <p class="font-bold">Note:</p>
        <p>
          Your account does not currently have MFA set up. Enabling MFA adds an
          additional layer of security to protect your account from unauthorized
          access.
        </p>
      </div>
    </div>
    
    <div class="flex w-1/2 flex-col items-center justify-center">
      <qrcode-vue
        :value="url"
        :level="level"
        :render-as="renderAs"
        :size="200"
      />
      <p class="mt-4 text-center">
        Can't scan the QR code? No problem! You can manually enter the MFA Key
        in your authenticator app:
      </p>
      <p class="mt-2 text-lg font-bold">{{ secret }}</p>

      <action-button
        class="mb-4 mt-4 rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
        text="Configure"
        @click="handleConfigureMFA(email, secret)"
      />

      <p class="mt-4 text-center text-sm italic">
        Click Configure only after setting up MFA with Google Authenticator!
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import QrcodeVue from "qrcode.vue";
import type { Level, RenderAs } from "qrcode.vue";
import ActionButton from "@/components/ActionButton.vue";
import type { CONFIGURE_MFA_REQUEST } from "@/types/user";
import { CONFIGURE_MFA } from "@/api/user";
import { AxiosError } from "axios";

const level = ref<Level>("M");
const renderAs = ref<RenderAs>("svg");

const responseMessage = ref("");

defineProps({
  email: {
    type: String,
    required: true,
  },
  secret: {
    type: String,
    required: true,
  },
  url: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["mfaConfigured"]);

const handleConfigureMFA = async (email: string, secret: string) => {
  try {
    const body: CONFIGURE_MFA_REQUEST = {
      email,
      secret,
    };

    const { status } = await CONFIGURE_MFA(body);

    if (status === 204) {
      emit("mfaConfigured", true);
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else {
      responseMessage.value = "Unexpected error occurred";
    }
  }
};
</script>
