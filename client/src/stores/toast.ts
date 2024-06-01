import { defineStore } from "pinia";
import { ref } from "vue";

export const useToastStore = defineStore("toast", () => {
  const toast = ref({
    visible: false,
    message: "",
    toastType: "",
    timer: ref<ReturnType<typeof setTimeout> | null>(null),
  });

  const SUCCESS_TOAST = (message: string, timeoutInSeconds: number = 2) => {
    SHOW_TOAST(message);

    toast.value.toastType = "success";

    toast.value.timer = setTimeout(() => {
      HIDE_TOAST();
    }, timeoutInSeconds * 1000);
  };

  const ERROR_TOAST = (message: string, timeoutInSeconds: number = 2) => {
    SHOW_TOAST(message);

    toast.value.toastType = "error";

    toast.value.timer = setTimeout(() => {
      HIDE_TOAST();
    }, timeoutInSeconds * 1000);
  };

  const SHOW_TOAST = (message: string) => {
    toast.value.visible = true;
    toast.value.message = message;
  };

  const HIDE_TOAST = () => {
    toast.value.visible = false;
    clearTimeout(toast.value.timer!);
  };

  return {
    toast,
    SUCCESS_TOAST,
    ERROR_TOAST,
  };
});
