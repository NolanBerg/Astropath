package core

import (
	"io"
	"os"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err) && !info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	dir, err := os.Open(path)
	if err != nil {
		return false
	}

	defer dir.Close()
	return true
}

// IsEmptyDirectory checks if the given path is an existing directory and is empty.
func IsEmptyDirectory(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	dir, err := os.Open(path)
	if err != nil {
		return false
	}

	defer dir.Close()

	entries, err := dir.ReadDir(1)
	if err != nil && err != io.EOF {
		return false
	}

	return len(entries) == 0
}
