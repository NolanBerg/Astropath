// All typescript types can go here.

export enum BlendingMode {
  Brighten = 0,
  Darken = 1,
}

export enum OutputFormat {
  TIFF = 0,
  JPEG = 1,
}

export interface Workflow {
  UID: number;
  Name: string;
  BlendingMode: BlendingMode;

  OutputLocation: string;
  OutputFileName: string;
  OutputFormat: OutputFormat;

  CreateTimelapseVideo: boolean;
  TimelapseLocation: string;
  TimelapseFramesLocation: string; // Takes up a massive amount of disk space as each frame is coppied 2x on disk. Named Sequentially (1.TIFF | 1.JPEG ... -> n.JPEG)
  DeleteFramesAfterProcessing: boolean;
  TimelapseDuration: number;
}

export interface DoctorReport {
  SystemHasFFMPEG: boolean;
  SystemHasARWConversion: boolean;
}

export interface ProcessStatus {
  Processed: number;
  Total: number;
  PreviewPath: string;
}

export interface AppSettings {
  UserFirstTime: boolean; // whether this is the user's first time experiencing the application or not. set to true after leaving the main-menu page for first time
  Workflows: Workflow[]; // Slice of all Workflow's the user has saved.

  EnableTimelapseGeneration: boolean; // Whether timelapse generation support is allowed. If true, user MUST have ffmpeg
  EnableARWConversion: boolean; // Whether arw conversion is enabled. If true the user must have dcraw installed
  ARWTempFilePath: string; // Where to store the files for conversion. When the user enables this in the application, the user MUST enter a valid location for this. It will take up alot of space.
}

export interface ImportResult {
  ErrorMessage: string;
  FilePaths: string[];
  ImageBounds: any;
}

export interface ImageBatch {
  FilePaths: string[];
  ImageBounds: any;
}
