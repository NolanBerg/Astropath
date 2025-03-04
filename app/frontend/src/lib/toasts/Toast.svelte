<script lang="ts">
    import toastManager, { type Toast } from "./toast";
    export let toast: Toast;
</script>

<div class={`alert alert-${toast.kind} alert-dismissible fade show`} role="alert">
  <strong>{(toast.kind == "warning" || toast.kind == "danger" )? "Error" : "Info"}</strong> {toast.message}
  <!-- <button type="button" class="close" data-dismiss="alert" aria-label="Close">
    <span aria-hidden="true" on:click={() => toastManager.removeToast(toast.uid)}>&times;</span>
  </button> -->
</div>

