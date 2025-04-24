import { writable } from "svelte/store";

export type ToastKind = "warning" | "info" | "danger" | "primary";
export interface Toast {
  uid: string;
  kind: ToastKind;
  message: string;
}

export function createToastStore() {
  const storeValue = writable<Toast[]>([]);

  function addToast(toast: Toast) {
    storeValue.update((prev) => [...prev, toast]);
  }

  function removeToast(uid: string) {
    storeValue.update((prev) => prev.filter((t) => t.uid !== uid));
  }

  function removeToastAfterTime(uid: string, ms: number) {
    setTimeout(() => removeToast(uid), ms);
  }

  return {
    set: storeValue.set,
    subscribe: storeValue.subscribe,
    update: storeValue.update,
    addToast,
    removeToast,
    removeToastAfterTime,
  };
}

const toastManager = createToastStore();
export default toastManager;

export function displayToast(
  message: string,
  kind: ToastKind,
  duration = 6000,
) {
  const toast = {
    message,
    kind,
    uid: crypto.randomUUID(),
  };

  toastManager.addToast(toast);
  toastManager.removeToastAfterTime(toast.uid, duration);
}
