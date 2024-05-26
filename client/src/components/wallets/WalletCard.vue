<template>
  <div :class="cardClass">
    <h2 class="mb-2 text-xl font-bold">{{ formattedType }}</h2>
    <div class="mb-4 flex gap-2 text-lg">
      <div>{{ currencyAmount }}</div>
    </div>
    <div class="flex items-center justify-between">
      <router-link
        class="text-blue-500 underline underline-offset-4"
        :to="{ name: '' }"
        >History</router-link
      >
      <action-button
        class="hover rounded-lg bg-blue-500 px-4 py-2 text-white transition hover:bg-blue-300"
        text="Transfer"
      />
    </div>
  </div>
</template>
<script lang="ts" setup>
import ActionButton from "@/components/ActionButton.vue";
import { computed } from "vue";

const props = defineProps({
  type: {
    type: String,
    required: true,
  },
  currencyAmount: {
    type: Array,
    required: true,
  },
});

const formattedType = computed(() => {
  return props.type.charAt(0).toUpperCase() + props.type.slice(1);
});

const cardClass = computed(() => {
  let dynamicClass = "rounded-lg border p-6 shadow-lg dark:text-black";

  const typeClassMapping: Record<string, string> = {
    personal: "bg-purple-300 shadow-lg",
    savings: "bg-green-300 shadow-lg",
    investment: "bg-blue-300 shadow-lg",
    business: "bg-red-300 shadow-lg",
  };

  if (props.type in typeClassMapping) {
    dynamicClass += " " + typeClassMapping[props.type];
  }

  return dynamicClass;
});
</script>
