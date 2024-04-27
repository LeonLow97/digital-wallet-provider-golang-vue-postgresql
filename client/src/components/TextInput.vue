<template>
  <input
    type="text"
    class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-white"
    :value="modelValue"
    @input="handleInput"
  />
</template>

<script lang="ts" setup>
defineProps({
  modelValue: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:modelValue"]);

const handleInput = ($event: Event) => {
  const target = $event.target as HTMLInputElement;
  emit("update:modelValue", target.value);
};
</script>

<style scoped></style>
