<script lang="ts">
    import { link } from "svelte-spa-router";
    import type { ProcessStatus, Workflow } from "../../../src/models/types";
    import { go } from "../../api";

    export let wf: Workflow;

    const startTime = performance.now();
    let elapsed: number;

    let status: ProcessStatus = {
        Processed: 0,
	    Total: 0,
	    PreviewPath: ""
    };

    let generatingTimelapse = false;
    $: finished = status.Total === status.Processed && status.Processed !== 0 && !generatingTimelapse;
    $: progress = status.Total > 0 ? (status.Processed / status.Total) * 100 : 0;

    window["runtime"].EventsOn('frame:finish', (s: ProcessStatus) => {
        status = s;
        elapsed = (performance.now() - startTime) / 1000;
    });

    window["runtime"].EventsOn('workflow:timelapse-generation-start', () => {
        generatingTimelapse = true;
    });

    window["runtime"].EventsOn('workflow:timelapse-generation-finished', () => {
        generatingTimelapse = false;
    });
</script>

<div class="container mt-4">
    <div class="card">
        <div class="card-header">
            {generatingTimelapse? "Generating Timelapse": "Processing Images"} [{status.Processed}/{status.Total}]
        </div>
        <div class="card-body">
            <div class="progress mb-3">
                <div 
                    class="progress-bar" 
                    role="progressbar" 
                    style="width: {progress}%" 
                    aria-valuenow={progress} 
                    aria-valuemin="0" 
                    aria-valuemax="100">
                    {Math.round(progress)}%
                </div>
            </div>

            {#if status.PreviewPath}
                <div class="d-flex justify-content-center">
                    <img 
                        src={status.PreviewPath + "?r="+Date.now()} 
                        alt="Preview" 
                        class="img-fluid rounded"
                        style="max-height: 82vh; width: auto;"
                        class:imgfinished={finished}
                    />
                </div>
            {/if}
        </div>
    </div>

    {#if finished}
        <div class="finished container">
            <h2>Processing Finished</h2>
            <small>Completed processing of {status.Total} images in {elapsed.toFixed(1)}s</small>
            <hr>

            <div class="btns-row mt-4">
                <a class="btn btn-primary" href="/main-menu" use:link>Main Menu</a>
                <button class="btn btn-outline-dark" on:click={() => go.ShowPathInFinder(status.PreviewPath)}>View in Finder</button>
                
                {#if wf.CreateTimelapseVideo}
                    {#if !wf.DeleteFramesAfterProcessing}
                        <button class="btn btn-outline-dark" on:click={() => go.ShowPathInFinder(wf.TimelapseFramesLocation)}>View Timelapse Frames in Finder</button>
                    {/if}
                    
                    <button class="btn btn-outline-dark" on:click={() => go.ShowPathInFinder(wf.TimelapseLocation)}>View Timelapse In Finder</button>
                {/if}
            </div>
        </div>
    {/if}
</div>


<style>
    .imgfinished {
        max-height: 65vh !important;
    }
</style>