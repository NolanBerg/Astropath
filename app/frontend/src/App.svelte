<script lang="ts">
    import Router, { push } from "svelte-spa-router";
    import type { DoctorReport } from "./models/types";
    import MainMenu from "./routes/main-menu/MainMenu.svelte";
    import Settings from "./routes/settings/Settings.svelte";
    import Upload from "./routes/upload/Upload.svelte";
    import Sequencing from "./routes/sequencing/Sequencing.svelte";
    import Workflow from "./routes/workflow/Workflow.svelte";

    window["runtime"].EventsOn('doctor:start', () => {
        console.log('Doctor started');
    });
    
    window["runtime"].EventsOn('doctor:complete', (report: DoctorReport) => {
        console.log('Doctor completed', report);
    });

    const routes = {
        "/main-menu": MainMenu,
        
        "/settings": Settings,
        
        "/upload": Upload,  // User will upload files and see the progress of the file checking. If any conversion occurs, they will see the progress here too
        
        "/sequencing": Sequencing, // User can select what images will be imported aswell as in what order.
        
        "/workflow/:uid": Workflow, // Workflow route is where the user selects the options for the processing. 
        // They can also select a diferent workflow, save, edit existing workflows. Accepts a Workflow UID as a parameter
        
        // "/processing:id": ProcessWorkflow, // will show the processing progress. When processing is done, show the final image/and or video.
        
        "*": MainMenu, // Fallback route
    };

</script>


<Router {routes} />