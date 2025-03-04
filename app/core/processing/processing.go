package processing

import (
	"fmt"
	"image"
	"path/filepath"

	"time"

	"app/core/events"
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
	// TODO! Check if the output folder exists. If it does not, create it
	// If it cannot be creaed, error out.
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
			// TODO: Instead of aborting, maybe skip???
			return err
		}
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

	outputFilePath := filepath.Join(wp.Workflow.OutputLocation, wp.Workflow.OutputFileName)
	fmt.Printf(" outputFilePath= %s\n", outputFilePath)
	imaging.Save(wp.ImageOut, outputFilePath, imaging.JPEGQuality(100)) // todo handle errors

	// Write file to output location
	wp.Processed++
	status := ProcessStatus{
		Processed:   int(wp.Processed),
		Total:       len(wp.Batch.FilePaths),
		PreviewPath: outputFilePath,
	}

	fmt.Printf("Processed %v (%d/%d)\n", incomingPath, wp.Processed, len(wp.Batch.FilePaths))
	wp.eventManager.Emit("frame:finish", status)

	// TODO: Cleanup the ARW TMP FILES FOLDER!!!!!!!

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
