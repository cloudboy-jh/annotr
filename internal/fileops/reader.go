package fileops

import (
	"os"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
