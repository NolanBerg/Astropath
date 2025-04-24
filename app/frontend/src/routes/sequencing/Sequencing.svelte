<script lang="ts">
    import { link } from "svelte-spa-router";
    import { batchStore } from '../../models/stores';
    import LazyImage from './LazyImage.svelte';

    $: filePaths = $batchStore.FilePaths;
    $: fileNames = filePaths.map(filePath => {
        const lastForward = filePath.lastIndexOf("/");
        const lastBack = filePath.lastIndexOf("\\");
        const lastSlash = Math.max(lastForward, lastBack);
        return lastSlash === -1 ? filePath : filePath.substring(lastSlash + 1);
    });

    let showFileNames = true;
    let showThumbnails = false;

    function isValidIndex(arr: string[], ind: number): boolean {
        return ind >= 0 && ind < arr.length;
    }

    function swapUp(arr: string[], ind: number): void {
        if (!isValidIndex(arr, ind) || ind === 0) return;
        [arr[ind], arr[ind - 1]] = [arr[ind - 1], arr[ind]];
        filePaths = filePaths;
    }

    function swapDown(arr: string[], ind: number): void {
        if (!isValidIndex(arr, ind) || ind === arr.length - 1) return;
        [arr[ind], arr[ind + 1]] = [arr[ind + 1], arr[ind]];
        filePaths = filePaths;
    }
</script>

<div class="container pe-4">
    <br>
    <h1 class="text-center me-4">Sequencing Handler</h1>
    <br><br>
</div>

<div class="container-fluid main-layout ps-4">
    <div class="scrollable-column ps-2">
        <h3>Files selected: {filePaths.length}</h3>
        <div class="image-list">
            {#each filePaths as path, index}
                <div class="card mb-2">
                    <div class="card-body d-flex justify-content-between align-items-center"> 
                           
                        <div class="d-flex flex-column m-1 swap">
                            <h3 class="img-num">{index + 1}</h3>
                            <button class="swap-btn" on:click={() => swapUp(filePaths, index)}>
                                <i class="fa-regular fa-square-caret-up swap-arrows"></i>
                            </button>
                            <button class="swap-btn" on:click={() => swapDown(filePaths, index)}>
                                <i class="fa-regular fa-square-caret-down swap-arrows"></i>
                            </button>
                        </div>
                        {#if showFileNames}
                            <p class="mb-0">{fileNames[index]}</p>
                        {:else}
                            <p class="mb-0"><strong>Image {index + 1}</strong></p>
                        {/if}
                        {#if showThumbnails}
                            <LazyImage src={path} width="300" alt="Thumbnail" />
                        {/if}
                        <button class="btn btn-danger" on:click={() => batchStore.removeFile(path)}>
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
            {/each}
        </div>
    </div>

    <div class="fixed-column">
        <div class="content-wrapper">
            <div class="toggles">
                <div class="form-check form-switch">
                    <input class="form-check-input" type="checkbox" role="switch" id="showFileNamesSwitch" bind:checked={showFileNames}>
                    <label class="form-check-label" for="showFileNamesSwitch">Show filenames</label>
                </div>
                <div class="form-check form-switch">
                    <input class="form-check-input" type="checkbox" role="switch" id="showThumbnailsSwitch" bind:checked={showThumbnails}>
                    <label class="form-check-label" for="showThumbnailsSwitch">Show Thumbnails</label>
                </div>
            </div>
            <div class="buttons">
                <a class="btn btn-primary main" href="/main-menu" use:link>Main Menu</a>
                <a class="btn btn-primary second" href="/workflow/0" use:link>Edit Workflow</a>
            </div>
        </div>
    </div>
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
    background-color: #7E1E5A;
    color: white;
}

.swap {
    width: 50px;
    background-color: white;
}

.swap-arrows {
    font-size: x-large;
}

.swap-btn {
    border: 0;
    background-color: white;
}

.main-layout {
    display: flex;
    position: relative;
}

.scrollable-column {
    width: 70%;
    max-height: 70vh;
    overflow-y: auto;
    padding-right: 15px;
}

.fixed-column {
    position: fixed;
    right: 0;
    top: 0;
    width: 30%;
    height: 100vh;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: white;
}

.content-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 90%;
    width: 100%;
}

.toggles {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.buttons {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.form-check.form-switch {
    display: flex;
    align-items: center;
    width: 200px;
    margin-bottom: 10px;
}

.form-check-input {
    margin-right: 10px;
}

.form-check-label {
    flex: 1;
    text-align: left;
}

.img-num{
    text-align: center;
}

@media (max-width: 767px) {
    .main-layout {
        flex-direction: column;
    }

    .scrollable-column {
        width: 100%;
    }

    .fixed-column {
        position: static;
        width: 100%;
        height: auto;
        padding: 20px 0;
    }

    .content-wrapper {
        height: auto;
        justify-content: center;
        gap: 20px;
    }
}
</style>
