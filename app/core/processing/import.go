package processing

import (
	"app/core"
	"app/core/doctor"
	"app/core/events"
	"app/core/settings"
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/h2non/filetype"
)

// Given a slice of paths (representing valid files. Check that all files are the same size.)
// Will Open each file and Return the dimension of each file (since they are the same OR whatever error occured).
func ValidateFrames(files []string, eventManager *events.EventManager) (image.Rectangle, []string, error) {
	if len(files) < 2 {
		return image.Rectangle{}, nil, errors.New("at least two files are required to compare dimensions. files list is too sjort")
	}

	var expectedSize image.Rectangle
	var haveSize bool
	var numberOfFilesNeedingConversion = 0

	for i, filePath := range files {
		eventManager.Emit("validate-frames:progress", i, len(files))

		file, err := os.Open(filePath)
		if err != nil {
			return image.Rectangle{}, nil, fmt.Errorf("failed to open the file %q: %w", filePath, err)
		}

		// decode image and check dimensions
		img, err := imaging.Decode(file)
		_ = file.Close()
		if err != nil {
			return image.Rectangle{}, nil, fmt.Errorf("failed to decode file %q: %w", filePath, err)
		}

		if !haveSize {
			expectedSize = img.Bounds()
			haveSize = true
		} else if img.Bounds() != expectedSize {
			return image.Rectangle{}, nil, fmt.Errorf("image %q has dimensions %v, but expected %v; cannot process images of different dimensions", filePath, img.Bounds(), expectedSize)
		}

		if IsARWImage(filePath) {
			numberOfFilesNeedingConversion++
		} else if IsBaseImage(filePath) {
			if numberOfFilesNeedingConversion > 0 {
				return image.Rectangle{}, nil, fmt.Errorf("files of diferent types were encountered. Cannot have .arw and JPEG/PNG in same batch import")
			}
		} else {
			return image.Rectangle{}, nil, fmt.Errorf("unsupported filetype at %s. Astropath only supports .png, .jpg, .tiff and .arw files", filePath)
		}
	}

	// ARW File Conversion
	if numberOfFilesNeedingConversion > 0 {
		var convertedPaths []string
		if numberOfFilesNeedingConversion != len(files) {
			return image.Rectangle{}, nil, fmt.Errorf("files of diferent types were encountered. Cannot have .arw and JPEG/PNG in same batch import")
		}

		if !doctor.GetDoctorReport(eventManager).SystemHasARWConversion {
			return image.Rectangle{}, nil, fmt.Errorf("arw files were found and need to be converted. however, arw conversion is not enabled in this application")
		}

		appSettings, err := settings.LoadAppSettings()
		if err != nil {
			return image.Rectangle{}, nil, err
		}

		if !appSettings.EnableARWConversion {
			return image.Rectangle{}, nil, fmt.Errorf("arw files were found and need to be converted. however, arw conversion is not enabled in this application")
		}

		os.RemoveAll(appSettings.ARWTempFilePath)
		os.MkdirAll(appSettings.ARWTempFilePath, 0777)

		if !core.IsDir(appSettings.ARWTempFilePath) {
			return image.Rectangle{}, nil, fmt.Errorf("%s is not a directory", appSettings.ARWTempFilePath)
		}

		ts := time.Now()
		for i, filePath := range files {
			eventManager.Emit("arw-conversion:progress", i, len(files))
			outputPath, err := doctor.ConvertARWToTIFF(filePath, appSettings.ARWTempFilePath)
			if err != nil {
				os.RemoveAll(appSettings.ARWTempFilePath)
				os.MkdirAll(appSettings.ARWTempFilePath, 0777)
				return image.Rectangle{}, nil, fmt.Errorf("something went wrong converting the .arw files to .tiff files")
			}

			convertedPaths = append(convertedPaths, outputPath)
			fmt.Printf("converted %s to %s\n", filePath, outputPath)
		}

		// last minute sanity check that all files were converted
		if len(convertedPaths) != len(files) {
			os.RemoveAll(appSettings.ARWTempFilePath)
			os.MkdirAll(appSettings.ARWTempFilePath, 0777)
			return image.Rectangle{}, nil, fmt.Errorf("something went wrong converting the .arw files to .tiff files")
		}

		fmt.Printf("\n-----\nfile converstion took %v\n-----\n", time.Since(ts))
		return expectedSize, convertedPaths, nil
	}

	return expectedSize, files, nil
}

// Detect if a file is .ARW. Should be run before IsBaseImage()
func IsARWImage(filepath string) bool {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return false
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		fmt.Println("Error matching file type: ", err)
		return false
	}

	if kind.Extension == "tif" && strings.HasSuffix(strings.ToLower(filepath), ".arw") {
		_, err := imaging.Open(filepath)
		if err == nil {
			return true
		}
	}

	return false
}

// Image which does not need conversion and works out of the box. (JPEG, PNG, TIFF)
func IsBaseImage(filepath string) bool {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return false
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return false
	}

	return kind.MIME.Value == "image/jpeg" || kind.MIME.Value == "image/png" || kind.MIME.Value == "image/tiff"
}
