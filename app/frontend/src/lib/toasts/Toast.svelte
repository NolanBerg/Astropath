<script lang="ts">
    import { type Toast } from "./toast";
    export let toast: Toast;
</script>

<div class={`alert alert-${toast.kind} alert-dismissible fade show`} role="alert">
  <strong>{(toast.kind == "warning" || toast.kind == "danger" )? "Error" : "Info"}</strong> {toast.message}
</div>

