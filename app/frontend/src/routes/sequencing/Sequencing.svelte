<script lang="ts">
    import { link } from "svelte-spa-router";
    import { batchStore } from '../../models/stores';

    $: filePaths = $batchStore.FilePaths;
</script>

<div class="container pe-4">
    <br>
    <h1 class="text-center me-4">Sequencing Handler</h1>
    <hr>

    <a class="btn btn-primary" href="/main-menu" use:link>Main Menu</a>
    <a class="btn btn-primary" href="/workflow/-1" use:link>Goto Default Workflow</a>

    <p class="mt-4">First design this route in Figma. Then once this route is nice and clean, then design the code. 
    We will need to make sure the import process is finished and the stores are setup before processing.</p>

    <p class="mt-4">I have added delete functionality aswell as this batch store so you can see how to properly do CRUD operations on a store. This should be a good asset. As always if you
        have any questions feel free to ask :).
    </p>
    <hr>
    
</div>


<div class="container col-md-8 mt-4">
    <h3>Fils selected: {filePaths.length}</h3>
    <div class="image-list">
        {#each filePaths as path, index} <!-- iterate over each entry from the store -->
            <div class="card mb-2">
                <div class="card-body d-flex justify-content-between align-items-center">
                    <p class="mb-0"><strong>{index}</strong> - {path}</p>
                    <button class="btn btn-danger" on:click={() => batchStore.removeFile(path)}> <!-- This is how you would delete a entry from the store using the method I crated -->
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        {/each}
    </div>
</div>

<br>