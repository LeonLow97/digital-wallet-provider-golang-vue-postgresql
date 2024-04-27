<template>
  <button class="border-none">{{ text }}</button>
</template>

<script lang="ts" setup>
defineProps({
  text: {
    type: String,
    required: true,
  },
});
</script>
