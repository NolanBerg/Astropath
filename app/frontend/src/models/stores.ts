import { get, writable } from "svelte/store";
import type { AppSettings, DoctorReport, ImageBatch, Workflow } from "./types";
import { go } from "../../src/api";

export const doctorReportStore = createDoctorReportStore();
export const settingsStore = createAppSettingsStore();
export const batchStore = createBatchStore();

function createAppSettingsStore() {
  const storeValue = writable<AppSettings>(null);

  /**
   * Load settings from disk. Makes call to go backend and then loads the value into the store.
   */
  async function loadSettings(): Promise<AppSettings> {
    const settings = await go.LoadApplicationSettings();
    storeValue.set(settings);
    return settings;
  }

  /**
   * Update settings in backend and refresh store
   */
  async function updateSettings(settings: AppSettings): Promise<void> {
    // @ts-ignore
    await go.UpdateSettings(settings);
    await loadSettings();
  }

  async function selectARWDirectory(): Promise<string> {
    try {
      const path = await go.SelectARWDirectory();
      return path;
    } catch (err) {
      console.error("Error selecting directory:", err);
      throw err;
    }
  }

  async function updateARWPath(newPath: string): Promise<void> {
    await go.UpdateARWTempFilePath(newPath);
    storeValue.update(settings => ({
      ...settings,
      ARWTempFilePath: newPath
    }));
  }
   
  async function getWorkflows(): Promise<Workflow[]> {
    const settings = await loadSettings();
    return settings.Workflows;
  }

  async function getWorkflow(uid: number): Promise<Workflow> {
    const workflows = await getWorkflows();
    return workflows.find((w) => w.UID === uid);
  }

  async function getDefaultWorkflow(): Promise<Workflow> {
    return await getWorkflow(-1);
  }

  async function resetSettings(): Promise<AppSettings> {
    const settings = await go.ResetAppSettings();
    storeValue.set(settings);
    return settings;
  }

  loadSettings();

  return {
    subscribe: storeValue.subscribe,
    set: storeValue.set,
    loadSettings,
    updateSettings,
    selectARWDirectory,
    updateARWPath,
    getWorkflows,
    getWorkflow,
    getDefaultWorkflow,
    resetSettings,
  };
}

function createDoctorReportStore() {
  const storeValue = writable<DoctorReport>(null);

  async function getDoctorReport(): Promise<DoctorReport> {
    const report = await go.GetDoctorResults();
    storeValue.set(report);
    return report;
  }

  function hasAllDependencies(): boolean {
    const value = get(storeValue);
    return value?.SystemHasARWConversion && value?.SystemHasFFMPEG;
  }

  getDoctorReport();

  return {
    subscribe: storeValue.subscribe,
    set: storeValue.set,
    getDoctorReport,
    hasAllDependencies,
  };
}

function createBatchStore() {
  const storeValue = writable<ImageBatch>({ FilePaths: [], ImageBounds: null });

  function addFiles(newFiles: string[]) {
    storeValue.update((batch) => {
      batch.FilePaths = [...batch.FilePaths, ...newFiles];
      return batch;
    });
  }

  function removeFile(fileToRemove: string) {
    storeValue.update((batch) => {
      batch.FilePaths = batch.FilePaths.filter((file) => file !== fileToRemove);
      return batch;
    });
  }

  function rearangeFiles(newOrder: string[]) {
    storeValue.update((batch) => {
      batch.FilePaths = newOrder;
      return batch;
    });
  }

  function getFiles(): string[] {
    return get(storeValue).FilePaths;
  }

  function getBatchValue(): ImageBatch {
    return get(storeValue);
  }

  return {
    subscribe: storeValue.subscribe,
    set: storeValue.set,
    addFiles,
    removeFile,
    rearangeFiles,
    getFiles,
    getBatchValue,
  };
}