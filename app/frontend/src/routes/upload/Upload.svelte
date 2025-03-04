<script lang="ts">
    import type { ImportResult } from "../../models/types";
    import { go } from "../../api";
    import { link, push } from "svelte-spa-router";
    import { batchStore } from "../../models/stores";
    import { displayToast } from "../../lib/toasts/toast";

    let importRes: Promise<ImportResult> = null;
    let processStatus = {
        Processed: 0,
        Total: 0,
    };

    let inValidationPhase = true;

    window["runtime"].EventsOn("validate-frames:progress", (validated: number, totalFrames: number) => {
            console.info(`validated ${validated} of ${totalFrames} frames!`);
            processStatus = {
                Processed: validated,
                Total: totalFrames,
            };

            inValidationPhase = true;
        }
    );

    window["runtime"].EventsOn("arw-conversion:progress", ( converted: number, totalFrames: number, filename: string, newPath: string) => {
            processStatus = {
                Processed: converted,
                Total: totalFrames,
            };

            inValidationPhase = false;
            console.info(`converting ${converted} of ${totalFrames} frames! Filename: ${filename} | Newpath: ${newPath}`);
        },
    );

    function importImages() {
        importRes = go.ImportImages();

        importRes.then((res) => {
            console.info(res);
            if (res.FilePaths?.length > 0) {
                batchStore.set({
                    FilePaths: res.FilePaths,
                    Bounds: res.Bounds,
                });

                return push("/sequencing");
            }

            displayToast(res.ErrorMessage, "danger");
            push("/main-menu");
        });
    }
</script>

<section class="container mt-4 pt-4 center-container">
    {#if importRes == null}
        <h2 class="mb-4">Please select your images</h2>
        <button class="btn btn-primary mb-3" on:click|once={importImages}
            >Import</button
        >
    {:else}
        {#await importRes}
            <div class="container mt-4">
                <div class="row justify-content-center">
                    <div class="col-12 col-md-8 col-lg-6">
                        <div class="card">
                            <div class="card-body">
                                <p class="card-text text-center">
                                    {inValidationPhase
                                        ? "Validating"
                                        : "Converting"} Images: {processStatus.Processed}
                                    / {processStatus.Total}
                                </p>
                                <div class="progress">
                                    <div
                                        class="progress-bar"
                                        role="progressbar"
                                        style="width: {processStatus.Processed / processStatus.Total * 100}%"
                                        aria-valuenow={processStatus.Processed}
                                        aria-valuemin="0"
                                        aria-valuemax={processStatus.Total}
                                    ></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        {/await}
    {/if}
</section>

<style>
.center-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100vh;
    width: 100vw;
}

button.btn,
a.btn {
    display: inline-block;
    width: 200px;
    text-align: center;
    padding: 10px;
    box-sizing: border-box;
}

button.btn {
    background-color: #1E5A7E;
}
a.btn {
    background-color: #7E1E5A ;
    color: white;
}
</style>

