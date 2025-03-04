<script lang="ts">
    import { link } from "svelte-spa-router";
    import { batchStore } from '../../models/stores';

    $: filePaths = $batchStore.FilePaths;
    let showFileNames = true;
    let showThumbnails = false;

    function isValidIndex(arr: string[], ind: number): boolean {
        return ind >= 0 && ind < arr.length;
    }

    function swapUp(arr: string[], ind: number): void {
        if (!isValidIndex(arr, ind) || ind === 0) return;
        [arr[ind], arr[ind - 1]] = [arr[ind - 1], arr[ind]];
        filePaths = filePaths; // Trigger reactivity
    }

    function swapDown(arr: string[], ind: number): void {
        if (!isValidIndex(arr, ind) || ind === arr.length - 1) return;
        [arr[ind], arr[ind + 1]] = [arr[ind + 1], arr[ind]];
        filePaths = filePaths; // Trigger reactivity
    }
</script>

<div class="container pe-4">
    <br>
    <h1 class="text-center me-4">Sequencing Handler</h1>
    <hr>
</div>

<div class="container row">
    <div class="col-md-8 scrollable-column">
        <h3>Fils selected: {filePaths.length}</h3>
    <div class="image-list">
        {#each filePaths as path, index} <!-- iterate over each entry from the store -->
            <div class="card mb-2">
                <div class="card-body d-flex justify-content-between align-items-center">    
                    <div class="d-flex flex-column m-1 swap">
                        <button class="swap-btn" on:click={() => swapUp(filePaths, index)}>
                            <i class="fa-regular fa-square-caret-up swap-arrows"></i>
                        </button>
                        <button class="swap-btn" on:click={() => swapDown(filePaths, index)}>
                            <i class="fa-regular fa-square-caret-down swap-arrows"></i>
                        </button>
                    </div>
                    {#if showFileNames}
                        <p class="mb-0"><strong>{index + 1}</strong> - {path}</p>
                    {:else}
                        <p class="mb-0"><strong>{index + 1}</strong></p>
                    {/if}
                   {#if showThumbnails}
                     <img src={path} width="300" alt="" />
                   {/if}
                    <button class="btn btn-danger" on:click={() => batchStore.removeFile(path)}> <!-- This is how you would delete a entry from the store using the method I crated -->
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        {/each}
    </div>
    </div>
    <div class="col-md-3 offset-md-1">
        <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckChecked" bind:checked={showFileNames}>
            <label class="form-check-label" for="flexSwitchCheckChecked">Show filenames</label>
        </div>
        <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckChecked" bind:checked={showThumbnails}>
            <label class="form-check-label" for="flexSwitchCheckChecked">Show Thumbnails</label>
        </div>
        <a class="btn btn-primary main" href="/main-menu" use:link>Main Menu</a>
        <a class="btn btn-primary second" href="/workflow/-1" use:link>Edit Workflow</a>
    </div>
</div>

<div class="container col-md-8 mt-4">
    
</div>

<br>

<style>
a.btn {
    display: inline-block;
    width: 200px;
    text-align: center;
    padding: 10px;
    box-sizing: border-box;
    margin-top: 8px;
}

.main {
    background-color: #1E5A7E;
}
.second {
    background-color: #7E1E5A ;
    color: white;
}
.swap{
    width: 25px;
}
.swap-arrows{
    font-size: large;
}
.swap-btn{
    border: 0;
}
</style>
