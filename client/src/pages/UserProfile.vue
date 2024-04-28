<template>
  <h1 class="text-xl font-bold">User Profile</h1>
  <div class="flex w-1/2 flex-col gap-6">
    <div class="flex justify-between gap-6">
      <div class="flex w-full flex-col gap-2 text-sm">
        <label for="firstName">First Name:</label>
        <text-input v-model.trim="firstName" placeholder="First Name" />
      </div>
      <div class="flex w-full flex-col gap-2 text-sm">
        <label for="lastName">Last Name:</label>
        <text-input v-model.trim="lastName" placeholder="Last Name" />
      </div>
    </div>
    <div class="flex flex-col gap-2 text-sm">
      <label for="username">Username:</label>
      <text-input v-model.trim="username" placeholder="Username" />
    </div>
    <div class="flex flex-col gap-2 text-sm">
      <label for="mobileNumber">Mobile Number:</label>
      <text-input v-model.trim="mobileNumber" placeholder="Mobile Number" />
    </div>
    <div class="flex flex-col gap-2 text-sm">
      <label for="email">Email:</label>
      <text-input v-model.trim="email" placeholder="Email" />
    </div>
    <action-button
      class="mb-4 w-full rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
      text="Save Changes"
      @click="handleSaveChanges"
    />
  </div>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/stores/user";
import { onMounted, ref } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";
import type { User } from "@/types/user";
import { UPDATE_USER } from "@/api/user";
import type { UPDATE_USER_REQUEST } from "@/types/user";
import { AxiosError } from "axios";

const userStore = useUserStore();

// Data Fields
const firstName = ref("");
const lastName = ref("");
const username = ref("");
const mobileNumber = ref("");
const email = ref("");

onMounted(async () => {
  firstName.value = userStore.user.firstName;
  lastName.value = userStore.user.lastName;
  username.value = userStore.user.username;
  email.value = userStore.user.email;
  mobileNumber.value = userStore.user.mobileNumber;
});

const handleSaveChanges = async () => {
  try {
    const req: UPDATE_USER_REQUEST = {
      first_name: firstName.value === "" ? null : firstName.value,
      last_name: lastName.value === "" ? null : lastName.value,
      username: username.value,
      mobile_number: mobileNumber.value,
      email: email.value,
    };

    const { status } = await UPDATE_USER(req);

    if (status === 204) {
      const user: User = {
        firstName: firstName.value,
        lastName: lastName.value,
        username: username.value,
        mobileNumber: mobileNumber.value,
        email: email.value,
      };

      userStore.SAVE_USER(user);

      alert("Changes Saved Successfully!");
    }
  } catch (error: unknown) {
    if (error instanceof AxiosError) {
      if (error.response) {
        alert(error.response?.data?.message);
      }
    } else {
      console.error("Unexpected error", error);
    }
  }
};
</script>
