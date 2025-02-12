<script lang="ts">
    import type { ImportResult } from "../../models/types";
    import { go } from "../../api";
    import { link, push } from "svelte-spa-router";
    import { batchStore } from "../../models/stores";

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

                push("/sequencing");
            }
        });
    }
</script>

<section class="container mt-4 pt-4">
    {#if importRes == null}
        <button class="btn btn-primary" on:click|once={importImages}
            >Import</button
        >
        <a class="btn btn-outline-dark" href="/sequencing" use:link
            >Sequencing</a
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
        {:then results}
            <!-- If we make it to this block, and the .then() above does not move us to the sequencing route then there was an error :( -->
            <div class="container">
                <h2 class="test-center">Import Failed</h2>
                <hr />
                <div class="alert alert-warning me-4">
                    {results.ErrorMessage}
                </div>

                <a class="btn btn-dark" href="/main-menu" use:link>Main Menu</a>
                <button class="btn btn-outline-primary" on:click|once={importImages}>Try Again</button>
            </div>
        {/await}
    {/if}
</section>
