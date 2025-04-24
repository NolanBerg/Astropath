import { BlendingMode, OutputFormat } from "./models/types";

export function blendingModeToString(mode: BlendingMode): string {
  switch (mode) {
    case BlendingMode.Brighten:
      return "Brighten";
    case BlendingMode.Darken:
      return "Darken";
    default:
      return "Unknown";
  }
}

export function outputFormatToString(format: OutputFormat): string {
  switch (format) {
    case OutputFormat.TIFF:
      return "TIFF";
    case OutputFormat.JPEG:
      return "JPEG";
    default:
      return "Unknown";
  }
}
