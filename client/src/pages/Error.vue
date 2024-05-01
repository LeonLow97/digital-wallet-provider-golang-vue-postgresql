<template>
  <h1 class="flex items-center justify-center gap-1 text-xl font-bold">
    Oops... Something went wrong!
    <svg-icon type="mdi" :path="mdiEmoticonSadOutline" />
  </h1>
  <p class="flex justify-center">
    We encountered an unexpected error. Our team has been notified and is
    working to fix it. Please try again later.
  </p>
</template>

<script lang="ts" setup>
import { mdiEmoticonSadOutline } from "@mdi/js";
import SvgIcon from "@jamescoyle/vue-icon";
</script>
