<script lang="ts">
    import { link, push, location } from "svelte-spa-router";
    import { settingsStore, batchStore,  } from "../../models/stores";
    import { blendingModeToString, outputFormatToString } from "../../helpers";
    import { go } from "../../../src/api";
    import { get } from "svelte/store";
    import Processing from "./Processing.svelte";
    import { displayToast } from "../../lib/toasts/toast";
    import { type Workflow } from "../../models/types";

    export let params = {} as { uid: string };
    let uid = parseInt(params?.uid) ?? null;
    let processingWorkflow = false;
    let workflows = [] as Workflow[];
    let editingWorkflow = false;
    let workflowCopy: Workflow;
    let workflow: Workflow;

    $: changesMade = JSON.stringify(workflow) !== JSON.stringify(workflowCopy);
    $: isDefaultWorkflow = uid === 0;

    async function init () {
      if (uid == null || uid == undefined || typeof uid !== "number") {
        displayToast(`Something went wrong loading workflow with UID=${uid}`, "danger");
        return await push("/main-menu");
      }

      workflow = await settingsStore.getWorkflow(uid);
      workflows = await settingsStore.getWorkflows();
      workflowCopy = Object.assign({}, workflow);
      processingWorkflow = false;
    }

    init();

    function gotoWorkflow (newUID: number) {
      uid = newUID;
      push(`/workflow/${uid}`);
      init();
    }
    
    async function startProcessing() {
      const batch = get(batchStore);
      // @ts-expect-error
      go.StartProcessingWorkflow(await settingsStore.getWorkflow(uid), batch);
      processingWorkflow = true;
    }

    async function startEditMode() {
      workflows = await settingsStore.getWorkflows();
      workflow = await settingsStore.getWorkflow(uid);
      editingWorkflow = true;
      workflowCopy = Object.assign({}, workflow);
    }

    async function cancelChanges() {
      workflow = await settingsStore.getWorkflow(uid);
      workflowCopy = Object.assign({}, workflow);
      editingWorkflow = false;
    }

    async function saveChanges() {
      if (workflow.Name.length === 0) {
        return displayToast("Please provide a valid Workflow name", "warning");
      }

      if (isDefaultWorkflow) {
        return displayToast("Cannot edit the default Workflow", "warning");
      }

      if (workflow.TimelapseDuration < 5) {
        return displayToast("Timelapse duration cannot be less than 5s", "warning");
      }

      if (workflow.TimelapseDuration > 60) {
        return displayToast("Timelapse duration cannot be greater than 1 minute", "warning");
      } 

      try {
          const settings = await settingsStore.loadSettings();
          // @ts-ignore
          workflow.BlendingMode = parseInt(workflow.BlendingMode);
          // @ts-ignore
          workflow.OutputFormat = parseInt(workflow.OutputFormat);
          settings.Workflows = settings.Workflows.map(wf => wf.UID === uid? workflow : wf);

          await settingsStore.updateSettings(settings);
          workflowCopy = Object.assign({}, workflow);
          displayToast("Workflow changes saved", "primary");
        } catch (err) {
            console.error(err);
            displayToast("Could not save workflow changes", "danger");
        } finally {
          editingWorkflow = false;
          workflow = await settingsStore.getWorkflow(uid);
          workflows = await settingsStore.getWorkflows();
        }
    }

    async function createNewWorkflow () {
      try {
        const defaultWorkflow = workflows.find(w => w.UID === 0);
        const newWorkflow = Object.assign({}, defaultWorkflow);
        newWorkflow.UID = workflows.length;

        const settings = await settingsStore.loadSettings();
        settings.Workflows.push(newWorkflow);
        await settingsStore.updateSettings(settings);
        uid = newWorkflow.UID;

        workflowCopy = Object.assign({}, workflow);
        displayToast(`Created new Workflow[${uid}]`, "primary");
      } catch (err) {
        console.error(err);
        displayToast("Could not create new Workflow", "danger");
      } finally {
        editingWorkflow = false;
        workflow = await settingsStore.getWorkflow(uid);
        workflows = await settingsStore.getWorkflows();
      }
    }

    async function deleteWorkflow() {
      if (isDefaultWorkflow) {
        return displayToast("Cannot delete the default Workflow", "warning");
      }

      const settings = await settingsStore.loadSettings();
      settings.Workflows = settings.Workflows.filter(wf => wf.UID !== uid);
      await settingsStore.updateSettings(settings);
      displayToast(`Deleted workflow[${uid}] successfuly`, "info");
      gotoWorkflow(0);
    }

</script>

{#if processingWorkflow}
    <Processing wf={workflow}/>
{:else}
    {#if !workflow}
        <div class="container">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Processing Workflow ...</span> <!-- This should not be loading long. ~1s probably unless the user has a crapper machine -->
            </div>
        </div>
    {/if}

    {#if workflow}
    <section class="row g-0">
        <div class="d-flex flex-column flex-shrink-0 p-3 text-white bg-dark" style="width: 280px; min-height: 100vh;">
          <span class="fs-4 mb-3">Workflows
            {#if !editingWorkflow}
                <button class="ms-2 btn btn-sm btn-outline-dark" on:click={createNewWorkflow}><i class="fa fa-plus"></i></button>

                {#if !isDefaultWorkflow}
                  <button class="btn btn-sm btn-dark" on:click={startEditMode}><i class="fa fa-pencil"></i></button>
                  <button class="btn btn-sm btn-danger" on:click={deleteWorkflow}><i class="fa fa-trash"></i></button>
                {/if}
            {/if}
          </span>

          <hr>
          <ul class="nav nav-pills flex-column mb-auto">
            {#each workflows as wf}
            <li class="nav-item">
                <!-- svelte-ignore a11y-missing-attribute -->
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <a class="nav-link" on:click={() => gotoWorkflow(wf.UID)} class:disabled={editingWorkflow} class:active={wf.UID === workflow.UID}>
                    {wf.Name}   
                </a>
            </li>
            {/each}
          </ul>
        </div>
      
        <div class="col p-4">
            <h1 class="mb-4">{workflow.Name}</h1>

            <div class="mb-4">
              {#if !editingWorkflow}
                <button class="btn btn-dark me-4" style="max-width: 50%;" on:click={startProcessing}>Process Workflow</button>
              {:else}
                <button class="btn btn-success me-2" on:click={saveChanges} disabled={!changesMade}>Save Changes <i class="fa fa-save"></i></button>
                <button class="btn btn-danger" on:click={cancelChanges}>Cancel <i class="fa fa-cancel"></i></button>
              {/if}
            </div>
          
            <div class="card">
              <div class="card-body">
                <div class="row">
                  <div class="col-md-6">
                    <ul class="list-unstyled">
                      <li>
                        <strong>UID:</strong> {workflow.UID}
                      </li>

                      <li>
                        <strong>Worklow Name: </strong>
                        {#if editingWorkflow}
                          <input type="text" class="form-control mt-1" bind:value={workflow.Name} required minlength="2" maxlength="40" />
                        {:else}
                          {workflow.Name}
                        {/if}
                      </li>

                      <li>
                        <strong>Output File Name:</strong> 
                        {#if editingWorkflow}
                          <input type="text" class="form-control mt-1" bind:value={workflow.OutputFileName} />
                        {:else}
                          {workflow.OutputFileName}
                        {/if}
                      </li>
                      
                      <li>
                        <strong>Blending Mode:</strong> 
                        {#if editingWorkflow}
                          <select class="form-select mt-1" bind:value={workflow.BlendingMode}>
                            <option value="0">Brighten</option>
                            <option value="1">Darken</option>
                          </select>
                        {:else}
                          {blendingModeToString(workflow.BlendingMode)}
                        {/if}
                      </li>
                      <li>
                        <strong>Output Format:</strong> 
                        {#if editingWorkflow}
                          <select class="form-select mt-1" bind:value={workflow.OutputFormat}>
                            <option value="0">TIFF</option>
                            <option value="1">JPEG</option>
                          </select>
                        {:else}
                          {outputFormatToString(workflow.OutputFormat)}
                        {/if}
                      </li>
                    </ul>
                  </div>
                  <div class="col-md-12">
                    
                    <ul class="list-unstyled">
                      <li>
                        <strong>Output Location:</strong> 
                        {#if editingWorkflow}
                          <input type="text" class="form-control mt-1" bind:value={workflow.OutputLocation} />
                        {:else}
                          {workflow.OutputLocation}
                        {/if}
                      </li>
                      <li>
                        <strong>Create Timelapse:</strong> 
                        {#if editingWorkflow}
                          <div class="form-check mt-1">
                            <input type="checkbox" class="form-check-input" id="createTimelapse" bind:checked={workflow.CreateTimelapseVideo} />
                            <label class="form-check-label" for="createTimelapse">Create Timelapse Video</label>
                          </div>
                        {:else}
                          {workflow.CreateTimelapseVideo ? "yes" : "no"}
                        {/if}
                      </li>
                      {#if workflow.CreateTimelapseVideo}
                      <li>
                        <strong>Timelapse Frame Location:</strong> 
                        {#if editingWorkflow}
                          <input type="text" class="form-control mt-1" bind:value={workflow.TimelapseFramesLocation} />
                        {:else}
                          {workflow.TimelapseFramesLocation}
                        {/if}
                      </li>
                      <li>
                        <strong>Delete Frames After Processing:</strong> 
                        {#if editingWorkflow}
                          <div class="form-check mt-1">
                            <input type="checkbox" class="form-check-input" id="deleteFramesAfterProcessing" bind:checked={workflow.DeleteFramesAfterProcessing} />
                            <label class="form-check-label" for="deleteFramesAfterProcessing">Delete Frames: </label>
                          </div>
                        {:else}
                          {workflow.DeleteFramesAfterProcessing ? "yes" : "no"}
                        {/if}
                      </li>
                      <li>
                        <strong>Timelapse Duration:</strong> 
                        {#if editingWorkflow}
                          <input type="number" step="1" min="5" max="60" class="form-control mt-1" bind:value={workflow.TimelapseDuration} />
                        {:else}
                          {workflow.TimelapseDuration}s
                        {/if}
                      </li>
                      {/if}
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>

      </section>
    {/if}

{/if}

