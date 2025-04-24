package doctor

import (
	"app/core/events"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type DoctorReport struct {
	SystemHasFFMPEG        bool
	SystemHasARWConversion bool
}

// ----------------------------------------------------
// OSX Install: https://formulae.brew.sh/formula/dcraw
// dcraw https://github.com/ncruces/dcraw
// dcraw -T -6 -W {INPUT_FILE.ARW}
// ----------------------------------------------------

// Accepts a input .arw file as input and output folder and will use the dcraw program to convert and move the selected desitnation folder
func ConvertARWToTIFF(inputFile string, outputFolder string) (string, error) {
	cmd := exec.Command("dcraw", "-T", "-6", "-W", inputFile)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert %s to TIFF: %v", inputFile, err)
	}

	baseName := strings.Split(filepath.Base(inputFile), ".")[0]
	dcrawOutput := filepath.Join(filepath.Dir(inputFile), baseName+".tiff")

	outputFilename := baseName + ".tiff"
	outputFilePath := filepath.Join(outputFolder, outputFilename)

	// move file to destinatiuon folder
	moveCmd := exec.Command("mv", dcrawOutput, outputFilePath)
	err = moveCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to move tiff file to output folder: %v", err)
	}

	return outputFilePath, nil
}

// Check the user operating system and install for FFMPEG for a installation of ffmpeg
func checkUserInstallForFFMPEG(wg *sync.WaitGroup, result *bool) {
	defer wg.Done()

	// Run "ffmpeg -version" command to check installation
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()

	// If no error, FFmpeg is installed
	*result = err == nil
}

// Check the user operating system and install for .awr converter
func checkUserInstallForARWConverter(wg *sync.WaitGroup, result *bool) {
	defer wg.Done()

	// Look for "dcraw" in the system PATH
	_, err := exec.LookPath("dcraw")

	// If found, dcraw is installed
	*result = err == nil
}

func GetDoctorReport(eventManager *events.EventManager) DoctorReport {
	eventManager.Emit("doctor:start")
	var wg sync.WaitGroup
	var userHasFFMPEG bool
	var userHasARWConverCapabilities bool

	wg.Add(2)
	go checkUserInstallForFFMPEG(&wg, &userHasFFMPEG)
	go checkUserInstallForARWConverter(&wg, &userHasARWConverCapabilities)
	wg.Wait()

	report := DoctorReport{
		SystemHasFFMPEG:        userHasFFMPEG,
		SystemHasARWConversion: userHasARWConverCapabilities,
	}

	eventManager.Emit("doctor:complete", report)
	return report
}
