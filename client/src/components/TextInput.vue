<template>
  <input
    :type="props.textInputType"
    class="w-full rounded-lg border border-gray-300 px-4 py-2 focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-white"
    :value="modelValue"
    @input="handleInput"
  />
</template>

<script lang="ts" setup>
const props = defineProps({
  modelValue: {
    type: [String, Number],
    required: true,
  },
  textInputType: {
    type: String,
    required: false,
    default: "text",
  },
});

const emit = defineEmits(["update:modelValue"]);

const handleInput = ($event: Event) => {
  const target = $event.target as HTMLInputElement;
  emit("update:modelValue", target.value);
};
</script>

<style scoped></style>
