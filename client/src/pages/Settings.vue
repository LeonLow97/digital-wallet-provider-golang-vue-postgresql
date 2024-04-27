<template>
  <div class="flex flex-col gap-4 pt-6">
    <h1 class="text-xl font-bold">Settings</h1>
    <label class="flex cursor-pointer items-center">
      <input
        type="checkbox"
        class="peer sr-only"
        v-model="isModeChecked"
        @click="handleMode"
      />
      <div
        class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-0.5 after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:ring-4 peer-focus:ring-blue-300 rtl:peer-checked:after:-translate-x-full dark:border-gray-600 dark:bg-gray-700 dark:peer-focus:ring-blue-800"
      ></div>
      <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">{{
        mode
      }}</span>
    </label>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";

// Data Fields
const isModeChecked = ref(false);
const mode = ref("Light Mode");

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
</script>
