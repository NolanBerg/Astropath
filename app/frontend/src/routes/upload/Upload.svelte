<script lang="ts">
    import type { ImportResult } from "../../models/types";
    import { go } from "../../api";
    import { push } from "svelte-spa-router";
    import { batchStore, settingsStore } from "../../models/stores";
    import { displayToast } from "../../lib/toasts/toast";
    
    let importRes: Promise<ImportResult> = null;
    let processStatus = {
        Processed: 0,
        Total: 0
    };
    let inValidationPhase = true;
    
    window["runtime"].EventsOn("validate-frames:progress", (validated: number, totalFrames: number) => {
        console.info(`validated ${validated} of ${totalFrames} frames!`);
        processStatus = {
            Processed: validated,
            Total: totalFrames
        };
        inValidationPhase = true;
    });
    
    window["runtime"].EventsOn("arw-conversion:progress", (converted: number, totalFrames: number) => {
        processStatus = {
            Processed: converted,
            Total: totalFrames
        };
        inValidationPhase = false;
    });
    
    function importImages() {
        // @ts-ignore
        importRes = go.ImportImages();
        importRes.then((res) => {
            console.info(res);
            if (res.FilePaths?.length > 0) {
                batchStore.set({
                    FilePaths: res.FilePaths,
                    //@ts-ignore
                    Bounds: res.Bounds
                });
                return push("/sequencing");
            }
            displayToast(res.ErrorMessage, "danger");
            push("/main-menu");
        });
    }
    </script>
    
    <div class="container vh-100 d-flex align-items-center justify-content-center">
        <div class="col-md-6">
            {#if importRes == null}
                <div class="card shadow-sm">
                    <div class="card-body text-center p-5">
                        <i class="fa fa-upload fa-4x text-primary mb-4"></i>
                        <h2 class="card-title mb-3">Import Image Batch</h2>
                        <p class="text-muted mb-4">
                            Select a batch of images to process. 
                            Supported formats include {#if $settingsStore.EnableARWConversion}
                            <strong>ARW</strong>,
                            {/if} <strong>TIFF</strong>, <strong>JPEG</strong> and <strong>PNG</strong>.
                        </p>

                        <button 
                            class="btn btn-primary btn-lg" 
                            on:click|once={importImages}>
                            Select Images <i class="fa fa-files ml-2"></i>
                        </button>
                    </div>
                </div>
            {:else}
                {#await importRes}
                    <div class="card shadow-sm">
                        <div class="card-body">
                            <p class="card-text text-center">
                                {inValidationPhase ? "Validating" : "Converting ARW"} Images: 
                                {processStatus.Processed} / {processStatus.Total}
                            </p>
                            <div class="progress">
                                <div
                                    class="progress-bar progress-bar-striped progress-bar-animated"
                                    role="progressbar"
                                    style="width: {(processStatus.Processed / processStatus.Total) * 100}%"
                                    aria-valuenow={processStatus.Processed}
                                    aria-valuemin="0"
                                    aria-valuemax={processStatus.Total}
                                ></div>
                            </div>
                        </div>
                    </div>
                {/await}
            {/if}
        </div>
    </div>
    
<style>

.card {
    border-radius: 10px;
}

</style>