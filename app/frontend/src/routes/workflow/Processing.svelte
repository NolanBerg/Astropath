<script lang="ts">
    import { type ProcessStatus } from "../../../src/models/types";
    let status: ProcessStatus = {
        Processed: 0,
	    Total: 0,
	    PreviewPath: ""
    };

    window["runtime"].EventsOn('frame:finish', (s: ProcessStatus) => {
        status = s;
        console.log(status);
    });

    $: progress = status.Total > 0 ? (status.Processed / status.Total) * 100 : 0;
</script>

<div class="container mt-4">
    <div class="card">
        <div class="card-header">
            Processing Images
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
                        style="max-height: 500px; width: auto;"
                    />
                </div>
            {/if}
        </div>
    </div>
</div>

