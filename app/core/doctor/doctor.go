package doctor

import (
	"app/core/events"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// Accepts a filepath to .arw as input and return new filepath as a result or error if one occured.
// Use a CLI application to do conversion.
func ConvertARWToTIFF(inputFile string) (string, error) {
	// Ensure the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Generate output filename by replacing .ARW with .TIFF
	outputFile := filepath.Join(
		filepath.Dir(inputFile), // Keep the original directory
		filepath.Base(inputFile[:len(inputFile)-len(filepath.Ext(inputFile))])+".tiff", // Change extension
	)

	// Run dcraw to convert ARW to TIFF
	cmd := exec.Command("dcraw", "-T", "-6", "-W", inputFile)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert %s to TIFF: %v", inputFile, err)
	}

	// Verify that output file was actually created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return "", fmt.Errorf("conversion failed: output file not found: %s", outputFile)
	}

	return outputFile, nil
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
