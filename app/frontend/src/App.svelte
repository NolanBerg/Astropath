<script lang="ts">
    import Router, { push } from "svelte-spa-router";
    import MainMenu from "./routes/main-menu/MainMenu.svelte";
    import Settings from "./routes/settings/Settings.svelte";
    import Upload from "./routes/upload/Upload.svelte";
    import Sequencing from "./routes/sequencing/Sequencing.svelte";
    import Workflow from "./routes/workflow/Workflow.svelte";
    import toastManager from "./lib/toasts/toast";
    import Toast from "./lib/toasts/Toast.svelte";

    const routes = {
        "/main-menu": MainMenu,
        
        "/settings": Settings,
        
        "/upload": Upload,  // User will upload files and see the progress of the file checking. If any conversion occurs, they will see the progress here too
        
        "/sequencing": Sequencing, // User can select what images will be imported aswell as in what order.
        
        "/workflow/:uid": Workflow, // Workflow route is where the user selects the options for the processing. 
        // They can also select a diferent workflow, save, edit existing workflows. Accepts a Workflow UID as a parameter
        
        "*": MainMenu, // Fallback route
    };
</script>

<Router {routes} />

<div class="container toast-container">
    <div class="toast-data">
        {#each $toastManager as toast}
            <Toast {toast}/>
        {/each}
    </div>
</div>

<style>
    .toast-container {
      position: absolute;
      bottom: 1em;
      right: 0.5em;
    }

    .toast-data {
        position: relative;
    }
</style>