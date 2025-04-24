<script lang="ts">
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";
    import { settingsStore, doctorReportStore } from "../../models/stores";
    import type { AppSettings } from "../../models/types";
    import { displayToast } from "../../../src/lib/toasts/toast";

    $: isFFMPEGMissing = !$doctorReportStore?.SystemHasFFMPEG;
    $: isDcrawMissing = !$doctorReportStore?.SystemHasARWConversion;

    let localSettings: AppSettings = { ...$settingsStore };
    onMount(() => {
        resetLocalSettings();
    });

    function resetLocalSettings() {
        localSettings = { ...$settingsStore };
    }

    async function handleTimelapseToggle(
        e: Event & { currentTarget: HTMLInputElement },
    ) {
        const enabled = e.currentTarget.checked;

        if (enabled && isFFMPEGMissing) {
            displayToast(
                "FFMPEG not found! Install FFMPEG to enable timelapse generation.",
                "warning",
            );
            e.currentTarget.checked = false;
            return;
        }

        localSettings.EnableTimelapseGeneration = enabled;
    }

    async function handleARWToggle(
        e: Event & { currentTarget: HTMLInputElement },
    ) {
        const enabled = e.currentTarget.checked;

        if (enabled && isDcrawMissing) {
            displayToast(
                "dcraw not found! Install dcraw to enable ARW conversion.",
                "danger",
            );
            e.currentTarget.checked = false;
            return;
        }

        localSettings.EnableARWConversion = enabled;
    }

    async function browseARWDirectory() {
        try {
            const path = await settingsStore.selectARWDirectory();
            if (path) {
                localSettings.ARWTempFilePath = path;
            }
        } catch (err) {
            displayToast(
                "Please choose a folder to store converted files.",
                "danger",
            );
        }
    }

    async function saveAndQuit() {
        try {
            if (
                localSettings.EnableARWConversion &&
                !localSettings.ARWTempFilePath
            ) {
                return displayToast(
                    "Failed to save settings: please select an output directory for ARW conversions",
                    "danger",
                );
            }

            await settingsStore.updateSettings(localSettings);
            displayToast("Settings saved successfully!", "primary");
            push("/main-menu");
        } catch (err: any) {
            console.error(err);
            displayToast("Failed to save settings: " + err.message, "danger");
        }
    }

    function quitWithoutSaving() {
        push("/main-menu");
    }

    async function resetSettings() {
        localSettings = await settingsStore.resetSettings();
        displayToast("Settings reset to defaults!", "primary");
    }
</script>

<div class="container height clean">
    <div class="row height">
        <div
            class="col-md-3 border-right height pt-4 accent-color d-flex flex-column align-items-center"
        >
            <button
                on:click={resetSettings}
                class="btn btn-secondary btn-block mint-color my-buttons"
            >
                Reset Settings
            </button>
            <button
                on:click={quitWithoutSaving}
                class="btn btn-secondary btn-block main-color my-buttons mt-2"
            >
                Quit
            </button>
            <button
                on:click={saveAndQuit}
                class="btn btn-secondary btn-block secondary-color my-buttons mt-2"
            >
                Save and Quit
            </button>
        </div>

        <!-- Main Content -->
        <div class="col-md-9 p-3">
            <h2>Astropath Settings</h2>

            <!-- Timelapse Generation Toggle -->
            <div class="form-check form-switch">
                <input
                    class="form-check-input"
                    type="checkbox"
                    role="switch"
                    bind:checked={localSettings.EnableTimelapseGeneration}
                    on:change={handleTimelapseToggle}
                    id="timelapseToggle"
                    disabled={isFFMPEGMissing}
                />
                <label class="form-check-label" for="timelapseToggle">
                    Enable Time Lapse Generation
                </label>
                {#if isFFMPEGMissing}
                    <p class="text-danger">
                        FFMPEG is not installed. Install FFMPEG to enable this
                        feature.
                    </p>
                {:else}
                    <p class="description">
                        Requires FFMPEG installation. Generates video timelapses
                        from image sequences.
                    </p>
                {/if}
            </div>

            <!-- ARW Conversion Toggle -->
            <div class="form-check form-switch">
                <input
                    class="form-check-input"
                    type="checkbox"
                    role="switch"
                    bind:checked={localSettings.EnableARWConversion}
                    on:change={handleARWToggle}
                    id="arwToggle"
                    disabled={isDcrawMissing}
                />
                <label class="form-check-label" for="arwToggle">
                    Support ARW Conversions
                </label>
                {#if isDcrawMissing}
                    <p class="text-danger">
                        dcraw is not installed. Install dcraw to enable this
                        feature.
                    </p>
                {:else}
                    <p class="description">
                        Requires dcraw installation. Allows conversion of Sony
                        RAW files to TIFF.
                    </p>
                {/if}
            </div>

            <!-- ARW Conversion Directory Selection -->
            {#if localSettings.EnableARWConversion}
                <div class="mt-3">
                    <label for="arwPath" class="form-label"
                        >ARW Conversion Directory</label
                    >
                    <div class="input-group">
                        <input
                            id="arwPath"
                            type="text"
                            bind:value={localSettings.ARWTempFilePath}
                            disabled
                            class="form-control"
                            placeholder="No directory selected"
                        />
                        <button
                            class="btn btn-outline-secondary"
                            type="button"
                            on:click={browseARWDirectory}
                        >
                            Browse
                        </button>
                    </div>
                    <p class="description mt-1">
                        Converted TIFF files will be stored in this directory.
                        After processing, files will be automatically deleted.
                    </p>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .border-right {
        border-right: 1px solid rgb(177, 177, 177);
    }
    .height {
        height: 100%;
    }
    .main-color {
        background-color: #1e5a7e;
    }
    .secondary-color {
        background-color: #7e1e5a;
    }
    .accent-color {
        background-color: #f0ede5;
    }
    .my-buttons {
        width: 80%;
    }
    .clean {
        margin-left: 0%;
        padding-left: 0;
    }
    .description {
        font-size: smaller;
        color: #8f8f8f;
    }
    .text-danger {
        font-size: smaller;
        color: red;
    }
    .mint-color {
        background-color: #4fd97d;
    }
    .input-group {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.5rem;
    }
    .input-group .form-control {
        flex: 1;
    }
</style>
