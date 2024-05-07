<template>
  <h1 class="text-xl font-bold">Settings</h1>

  <div>
    <div class="mb-4 text-sm text-slate-600">Toggle Mode</div>
    <label class="flex cursor-pointer items-center">
      <input
        type="checkbox"
        class="peer sr-only"
        v-model="isModeChecked"
        @click="handleMode"
      />
      <div
        class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-0.5 after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:ring-4 peer-focus:ring-blue-300 dark:border-gray-600 dark:bg-gray-700 dark:peer-focus:ring-blue-800 rtl:peer-checked:after:-translate-x-full"
      ></div>
      <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">{{
        mode
      }}</span>
    </label>
  </div>

  <div class="mt-5">
    <div class="mb-2 text-sm text-slate-600">Security</div>
    <button
      @click="openChangePasswordModal"
      class="rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
    >
      Change Password
    </button>
  </div>

  <Modal v-if="isModalOpen">
    <form @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-6">
        <!-- Modal content -->
        <h1 class="text-xl font-bold dark:text-white">Change Password</h1>
        <text-input
          v-model.trim="currentPassword"
          placeholder="Current Password"
          type="password"
          autocomplete="off"
        />
        <text-input
          v-model.trim="newPassword"
          placeholder="New Password"
          type="password"
          autocomplete="off"
        />
        <text-input
          v-model.trim="confirmPassword"
          placeholder="Confirm Password"
          type="password"
          autocomplete="off"
        />
        <div class="flex justify-end gap-4">
          <action-button
            @click="closeChangePasswordModal"
            class="mb-4 inline-block rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
            text="Close"
          />
          <action-button
            class="mb-4 inline-block rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
            text="Submit"
          />
        </div>
      </div>
    </form>
  </Modal>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import Modal from "@/components/Modal.vue";
import type { CHANGE_PASSWORD_REQUEST } from "@/types/user";
import { AxiosError } from "axios";
import { CHANGE_PASSWORD } from "@/api/user";

// Data Fields
const isModeChecked = ref(false);
const isModalOpen = ref(false);
const mode = ref("Light Mode");
const currentPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const responseMessage = ref("");

onMounted(async () => {
  if (localStorage.getItem("MODE") === "true") {
    isModeChecked.value = true;
    mode.value = "Dark Mode";
  }
});

const handleMode = () => {
  isModeChecked.value = !isModeChecked.value;
  mode.value = isModeChecked.value ? "Dark Mode" : "Light Mode";

  // get the html root element
  const htmlElement = document.documentElement;

  if (isModeChecked.value) {
    htmlElement.classList.add("dark");
    localStorage.setItem("MODE", "true");
  } else {
    htmlElement.classList.remove("dark");
    localStorage.setItem("MODE", "false");
  }
};

const openChangePasswordModal = () => {
  isModalOpen.value = true;
};
const closeChangePasswordModal = () => {
  isModalOpen.value = false;
};

const handleSubmit = async () => {
  // check if new password is same as confirm password
  if (newPassword.value !== confirmPassword.value) {
    alert("new password and confirm password do not match");
    return;
  }

  try {
    const body: CHANGE_PASSWORD_REQUEST = {
      current_password: currentPassword.value,
      new_password: confirmPassword.value,
    };

    console.log(body);

    const { status } = await CHANGE_PASSWORD(body);

    if (status === 204) {
      // reset data fields
      currentPassword.value = "";
      newPassword.value = "";
      confirmPassword.value = "";
      isModal.value = false; // close modal
      
      alert("Password changed successfully!");
    }
  } catch (error: any) {
    console.log(error);
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
