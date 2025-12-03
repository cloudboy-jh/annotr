package fileops

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudboy-jh/annotr/internal/parser"
)

type FileInfo struct {
	Path     string
	Name     string
	Language string
}

func ScanDirectory(dir string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		if parser.IsSupportedFile(path) {
			files = append(files, FileInfo{
				Path:     path,
				Name:     info.Name(),
				Language: getLanguageFromExt(filepath.Ext(path)),
			})
		}

		return nil
	})

	return files, err
}

func getLanguageFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	default:
		return ""
	}
}
