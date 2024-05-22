<template>
  <div>
    <!-- Overlay -->
    <div
      @click.self="closeModal"
      class="fixed left-0 top-0 flex h-full w-full items-center justify-center bg-gray-900 bg-opacity-50"
    >
      <!-- Modal -->
      <div @click.stop :class="modalCss">
        <!-- Modal content -->
        <div>
          <slot></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import ActionButton from "./ActionButton.vue";

const props = defineProps({
  modalWidth: {
    type: String,
    required: false,
    default: "3/5",
  },
});

const modalCss = computed(() => {
  return `w-${props.modalWidth} rounded-lg bg-white p-8 dark:bg-slate-500`;
});

const emit = defineEmits(["closeOverlay"]);

const closeModal = () => {
  emit("closeOverlay", true);
};
</script>
