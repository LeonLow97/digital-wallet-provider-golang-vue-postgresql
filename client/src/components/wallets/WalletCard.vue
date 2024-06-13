<template>
  <div :class="cardClass">
    <h2 class="mb-2 text-xl font-bold tracking-wider">{{ formattedType }}</h2>
    <div
      class="flex gap-2 text-lg tracking-wide"
      v-for="item in currencyAmount"
    >
      <div>{{ item.amount }} {{ item.currency }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import type { PropType } from "vue";
import type { WalletCurrencyAmount } from "@/types/wallet";

const props = defineProps({
  type: {
    type: String,
    required: true,
  },
  currencyAmount: {
    type: Array as PropType<WalletCurrencyAmount[]>,
    required: true,
  },
});

const formattedType = computed(() => {
  return props.type.charAt(0).toUpperCase() + props.type.slice(1);
});

const cardClass = computed(() => {
  let dynamicClass =
    "cursor-pointer rounded-lg border border-none p-6 shadow-lg dark:text-black";

  const typeClassMapping: Record<string, string> = {
    personal: "shadow-lg dark:bg-slate-700 dark:text-slate-300",
    savings: "shadow-lg dark:bg-slate-700 dark:text-slate-300",
    investment: "shadow-lg dark:bg-slate-700 dark:text-slate-300",
    business: "shadow-lg dark:bg-slate-700 dark:text-slate-300",
  };

  if (props.type in typeClassMapping) {
    dynamicClass += " " + typeClassMapping[props.type];
  }

  return dynamicClass;
});
</script>
