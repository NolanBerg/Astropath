package processing

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"time"

	"app/core/events"
	"app/core/settings"
	"app/core/workflow"

	"github.com/disintegration/imaging"
)

type ProcessStatus struct {
	Processed   int
	Total       int
	PreviewPath string
}

type ImageBatch struct {
	FilePaths   []string
	ImageBounds image.Rectangle
}

type WorkflowProcessor struct {
	StartTime    time.Time
	Processing   bool
	Processed    uint
	ImageOut     *image.NRGBA
	Batch        *ImageBatch
	Workflow     *workflow.Workflow
	eventManager *events.EventManager
}

func (wp *WorkflowProcessor) Process() error {
	// Perform cleanup of arw files, and timelapse frames at fn exit
	defer func() {
		if wp.Workflow.CreateTimelapseVideo && wp.Workflow.DeleteFramesAfterProcessing {
			fmt.Printf("Deleting frames in timelapse folder: %s\n", wp.Workflow.TimelapseFramesLocation)
			os.RemoveAll(wp.Workflow.TimelapseFramesLocation)
		}

		appSettings, _ := settings.LoadAppSettings()
		os.RemoveAll(appSettings.ARWTempFilePath)
		os.MkdirAll(appSettings.ARWTempFilePath, 0777)
		fmt.Printf("removing folder: %s\n", appSettings.ARWTempFilePath)
	}()

	if os.MkdirAll(wp.Workflow.OutputLocation, 0777) != nil {
		return fmt.Errorf("could not create output folder %s", wp.Workflow.OutputLocation)
	}

	if wp.Workflow.CreateTimelapseVideo {
		if os.MkdirAll(wp.Workflow.TimelapseFramesLocation, 0777) != nil {
			return fmt.Errorf("could not create timelapse frames folder %s", wp.Workflow.TimelapseFramesLocation)
		}

		if os.MkdirAll(wp.Workflow.TimelapseLocation, 0777) != nil {
			return fmt.Errorf("could not create output folder %s", wp.Workflow.TimelapseLocation)
		}
	}

	filesLen := len(wp.Batch.FilePaths)
	wp.eventManager.Emit("workflow:start")

	// if workflow not started yet, initialize its
	if !wp.Processing && wp.Processed == 0 {
		wp.StartTime = time.Now()
		wp.Processing = true
	}

	for wp.Processed < uint(filesLen) {
		err := wp.ProcessImage()
		if err != nil {
			fmt.Printf("[ERROR] Processing failed on %s: %v\n", wp.Batch.FilePaths[wp.Processed], err)
			return err
		}
	}

	// Create timelapse now that all the frames are created AND in the appropriate folder
	if wp.Workflow.CreateTimelapseVideo {
		wp.eventManager.Emit("workflow:timelapse-generation-start")
		outputVideo := filepath.Join(wp.Workflow.TimelapseLocation, "timelapse.mp4")

		files, err := os.ReadDir(wp.Workflow.TimelapseFramesLocation)
		if err != nil {
			return fmt.Errorf("failed to read frames directory: %v", err)
		}

		var images []os.DirEntry
		VALID_IMAGE_EXTENSIONS := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".tiff": true,
			".tif":  true,
		}

		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if VALID_IMAGE_EXTENSIONS[ext] {
				images = append(images, file)
			}
		}

		sort.Slice(images, func(i, j int) bool {
			num1, err1 := strconv.Atoi(strings.TrimSuffix(images[i].Name(), filepath.Ext(images[i].Name())))
			num2, err2 := strconv.Atoi(strings.TrimSuffix(images[j].Name(), filepath.Ext(images[j].Name())))

			if err1 != nil || err2 != nil {
				return images[i].Name() < images[j].Name()
			}
			return num1 < num2
		})

		if len(images) == 0 {
			return fmt.Errorf("no valid image files found in %s", wp.Workflow.TimelapseFramesLocation)
		}

		listFile, err := os.CreateTemp("", "timelapse-*.txt")
		if err != nil {
			return fmt.Errorf("failed to generate timelapse video")
		}

		defer func() {
			listFile.Close()
			os.Remove(listFile.Name())
		}()

		for _, file := range images {
			framePath := filepath.Join(wp.Workflow.TimelapseFramesLocation, file.Name())
			// escape single quotes for ffmpeg
			escapedPath := strings.ReplaceAll(filepath.ToSlash(framePath), "'", "'\\''")
			if _, err := fmt.Fprintf(listFile, "file '%s'\nduration 0.1\n", escapedPath); err != nil {
				return fmt.Errorf("failed to generate timelapse video")
			}
		}

		lastFrame := filepath.Join(wp.Workflow.TimelapseFramesLocation, images[len(images)-1].Name())
		escapedLastFrame := strings.ReplaceAll(filepath.ToSlash(lastFrame), "'", "'\\''")
		if _, err := fmt.Fprintf(listFile, "file '%s'\n", escapedLastFrame); err != nil {
			return fmt.Errorf("failed to generate timelapse video")
		}

		if err := listFile.Close(); err != nil {
			return fmt.Errorf("failed to generate timelapse video")
		}

		framerate := float64(len(images)) / float64(wp.Workflow.TimelapseDuration)
		cmd := exec.Command("ffmpeg",
			"-f", "concat",
			"-safe", "0",
			"-i", listFile.Name(),
			"-r", fmt.Sprintf("%.2f", framerate),
			"-c:v", "libx264",
			"-preset", "slow",
			"-crf", "18",
			"-pix_fmt", "yuv420p",
			"-y",
			outputVideo,
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		fmt.Printf("running ffmpeg command: %v", cmd.Args)
		if err := cmd.Run(); err != nil {
			fmt.Printf("stdout:\n%s", stdout.String())
			fmt.Printf("stderr:\n%s", stderr.String())
			return fmt.Errorf("failed to generate timelapse video")
		}

		wp.eventManager.Emit("workflow:timelapse-generation-finished")
	}

	wp.Processing = false
	wp.eventManager.Emit("workflow:finish", time.Since(wp.StartTime))
	return nil
}

func (wp *WorkflowProcessor) ProcessImage() error {
	wp.eventManager.Emit("frame:start", wp.Batch.FilePaths[wp.Processed])
	if wp.Processed >= uint(len(wp.Batch.FilePaths)) {
		return fmt.Errorf("all files already processed")
	}

	incomingPath := wp.Batch.FilePaths[wp.Processed]
	incomingImage, err := imaging.Open(incomingPath)
	if err != nil {
		return err
	}

	if wp.Processed == 0 {
		wp.ImageOut = imaging.Clone(incomingImage)
	} else {
		LuminanceBlendFrame(wp.Workflow.BlendingMode == workflow.BM_BRIGHTEN, wp.ImageOut, incomingImage)
	}

	fileExt := ".jpeg"
	if wp.Workflow.OutputFormat == workflow.OF_TIFF {
		fileExt = ".tiff"
	}

	outputFilePath := filepath.Join(wp.Workflow.OutputLocation, wp.Workflow.OutputFileName+fileExt)
	fmt.Printf(" outputFilePath= %s\n", outputFilePath)
	imaging.Save(wp.ImageOut, outputFilePath, imaging.JPEGQuality(100))

	// Write file to output location
	wp.Processed++
	status := ProcessStatus{
		Processed:   int(wp.Processed),
		Total:       len(wp.Batch.FilePaths),
		PreviewPath: outputFilePath,
	}

	if wp.Workflow.CreateTimelapseVideo {
		timelapseFilePath := filepath.Join(wp.Workflow.TimelapseFramesLocation, fmt.Sprintf("%d%s", wp.Processed, fileExt))
		imaging.Save(wp.ImageOut, timelapseFilePath, imaging.JPEGQuality(100))
	}

	fmt.Printf("Processed %v (%d/%d)\n", incomingPath, wp.Processed, len(wp.Batch.FilePaths))
	wp.eventManager.Emit("frame:finish", status)

	return nil
}

func CreateWorkflowProcessor(batch *ImageBatch, eventManager *events.EventManager, workflow *workflow.Workflow) *WorkflowProcessor {
	return &WorkflowProcessor{
		StartTime:    time.Time{},
		Processing:   false,
		Processed:    0,
		ImageOut:     image.NewNRGBA(batch.ImageBounds),
		Batch:        batch,
		Workflow:     workflow,
		eventManager: eventManager,
	}
}
