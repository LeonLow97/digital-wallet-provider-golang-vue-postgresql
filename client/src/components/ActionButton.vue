<template>
  <button class="shadow-md transition">{{ text }}</button>
</template>

<script lang="ts" setup>
defineProps({
  text: {
    type: String,
    required: true,
  },
});
</script>
