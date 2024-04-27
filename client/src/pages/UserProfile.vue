<template>
  <div class="flex flex-col gap-4 pt-6">
    <h1 class="text-xl font-bold">User Profile</h1>
    <div class="flex w-1/2 flex-col gap-6">
      <div class="flex justify-between gap-6">
        <div class="flex w-full flex-col gap-2 text-sm">
          <label for="firstName">First Name:</label>
          <text-input
            class="rounded-lg border px-4 py-2"
            v-model.trim="firstName"
            placeholder="First Name"
          />
        </div>
        <div class="flex w-full flex-col gap-2 text-sm">
          <label for="lastName">Last Name:</label>
          <text-input
            class="rounded-lg border px-4 py-2"
            v-model.trim="lastName"
            placeholder="Last Name"
          />
        </div>
      </div>
      <div class="flex flex-col gap-2 text-sm">
        <label for="username">Username:</label>
        <text-input
          class="rounded-lg border px-4 py-2"
          v-model.trim="username"
          placeholder="Username"
        />
      </div>
      <div class="flex flex-col gap-2 text-sm">
        <label for="mobileNumber">Mobile Number:</label>
        <text-input
          class="rounded-lg border px-4 py-2"
          v-model.trim="mobileNumber"
          placeholder="Mobile Number"
        />
      </div>
      <div class="flex flex-col gap-2 text-sm">
        <label for="email">Email:</label>
        <text-input
          class="rounded-lg border px-4 py-2"
          v-model.trim="email"
          placeholder="Email"
        />
      </div>
      <action-button
        class="mb-4 w-full rounded-lg border bg-blue-500 px-4 py-2 text-center text-white transition hover:bg-blue-400"
        text="Save Changes"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import TextInput from "@/components/TextInput.vue";
import ActionButton from "@/components/ActionButton.vue";

// Data Fields
const firstName = ref("Jie Wei");
const lastName = ref("Low");
const username = ref("lowjiewei");
const mobileNumber = ref("+65 98765432");
const email = ref("lowjiewei@email.com");

// Fetch user profile from backend
</script>
