package workflow

import "fmt"

type OutputFormat int

const (
	OF_TIFF = iota
	OF_JPEG
)

func (w OutputFormat) String() string {
	return []string{"TIFF", "JPEG"}[w]
}

type BlendingMode int

const (
	BM_BRIGHTEN = iota
	BM_DARKEN
)

func (w BlendingMode) String() string {
	return []string{"Brighten", "Darken"}[w]
}

type Workflow struct {
	UID          int          // date-timestamp
	Name         string       // workflow name
	BlendingMode BlendingMode // defaults to lighten

	OutputLocation string       // default to tmp folder
	OutputFileName string       // default to input file name
	OutputFormat   OutputFormat // default to jpeg

	CreateTimelapseVideo        bool
	TimelapseLocation           string
	TimelapseFramesLocation     string // Takes up a massive amount of disk space as each frame is coppied 2x on disk. Named Sequentially (1.TIFF | 1.JPEG ... -> n.JPEG)
	DeleteFramesAfterProcessing bool
	TimelapseDuration           int // duration in seconds
}

func CreateDefaultWorkflow(output_location string) *Workflow {
	return &Workflow{
		UID:            0,
		Name:           "Default Startrail Workflow",
		BlendingMode:   BM_BRIGHTEN,
		OutputFormat:   OF_TIFF,
		OutputLocation: output_location,
		OutputFileName: "lighten-startrail",

		CreateTimelapseVideo:        false, // Instead of creating a single output frame, need to create a frame for each step
		TimelapseLocation:           fmt.Sprintf("%s/timelapse/", output_location),
		TimelapseFramesLocation:     fmt.Sprintf("%s/timelapse_frames/", output_location),
		DeleteFramesAfterProcessing: true,
		TimelapseDuration:           15,
	}
}
