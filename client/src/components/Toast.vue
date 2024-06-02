<template>
  <div
    class="absolute left-1/2 top-5 max-w-md -translate-x-1/2 rounded-xl border border-none shadow-lg drop-shadow-md dark:border-neutral-700 dark:bg-neutral-800"
    role="alert"
    :class="dynamicCss('bg')"
  >
    <div class="flex items-center justify-center gap-2 p-4">
      <svg-icon
        type="mdi"
        :path="getIconPath(props.toastType)"
        :class="dynamicCss('icon')"
      />
      <p
        class="text-sm capitalize tracking-wider dark:text-neutral-400"
        :class="dynamicCss('text')"
      >
        {{ props.message }}
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { mdiCheckCircleOutline, mdiCloseCircleOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";

const props = defineProps<{
  message: string;
  toastType: string;
}>();

function getIconPath(toastType: string): string {
  if (toastType?.toLowerCase() === "success") {
    return mdiCheckCircleOutline;
  }
  return mdiCloseCircleOutline;
}

const dynamicCss = (property: "bg" | "text" | "icon") => {
  const typeClassMapping: Record<
    string,
    Record<"bg" | "text" | "icon", string>
  > = {
    success: {
      bg: "bg-green-100",
      text: "text-green-800",
      icon: "text-green-800",
    },
    error: {
      bg: "bg-red-100",
      text: "text-red-800",
      icon: "text-red-800",
    },
  };

  if (props?.toastType?.toLowerCase() in typeClassMapping) {
    return typeClassMapping[props.toastType.toLowerCase()][property];
  } else {
    return property === "bg" ? "bg-blue-100" : "text-blue-800";
  }
};
</script>
