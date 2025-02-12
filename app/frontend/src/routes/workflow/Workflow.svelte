<script lang="ts">
    import { push } from "svelte-spa-router";
    import { settingsStore } from "../../models/stores";
    import { blendingModeToString, outputFormatToString } from "../../helpers";

    export let params = {} as { uid: string };
    const uid = parseInt(params?.uid) ?? null;

    if (!uid) {
        console.error("no uid was provided. This is a HUGE programmer issue"); 
        // TODO: we should never reach this block since there should always be a UID, but it should be a good idea to do this check in the future.
        push("/main-menu"); // Goto the main menu
    }

</script>

{#await settingsStore.getWorkflow(uid)}
    <div class="container">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span> <!-- This should not be loading long. ~1s probably unless the user has a crapper machine -->
        </div>
    </div>
    {:then workflow}
    <div class="container mt-4">
        <h1 class="text-center mb-4">{workflow.Name}</h1>
        <div class="card">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <h5 class="card-title">Basic Information</h5>
                        <ul class="list-unstyled">
                            <li><strong>UID:</strong> {workflow.UID}</li>
                            <li><strong>Output File Name:</strong> {workflow.OutputFileName}</li>
                            <li><strong>Blending Mode:</strong> {blendingModeToString(workflow.BlendingMode)}</li>
                            <li><strong>Output Format:</strong> {outputFormatToString(workflow.OutputFormat)}</li>
                        </ul>
                    </div>
                    <div class="col-md-6">
                        <h5 class="card-title">Output Settings</h5>
                        <ul class="list-unstyled">
                            <li><strong>Output Location:</strong> {workflow.OutputLocation}</li>
                            <li><strong>Create Timelapse:</strong> {workflow.CreateTimelapseVideo ? 'Yes' : 'No'}</li>
                            
                            {#if workflow.CreateTimelapseVideo}
                                <li><strong>Timelapse Location:</strong> {workflow.TimelapseLocation}</li>
                                <li><strong>Frames Location:</strong> {workflow.TimelapseFramesLocation}</li>
                            {/if}
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
{/await}
