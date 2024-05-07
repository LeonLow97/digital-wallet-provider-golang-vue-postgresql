<template>
  <div>
    <!-- Overlay -->
    <div
      @click.self="closeModal"
      class="fixed left-0 top-0 flex h-full w-full items-center justify-center bg-gray-900 bg-opacity-50"
    >
      <!-- Modal -->
      <div @click.stop class="w-3/5 rounded-lg bg-white p-8 dark:bg-slate-500">
        <!-- Modal content -->
        <div>
          <slot></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import ActionButton from "./ActionButton.vue";
</script>
