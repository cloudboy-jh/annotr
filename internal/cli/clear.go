package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudboy-jh/annotr/internal/fileops"
	"github.com/cloudboy-jh/annotr/internal/parser"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear [file or directory]",
	Short: "Remove comments from files",
	Long: `Remove comments from a file or all supported files in a directory.

Examples:
  annotr clear file.go       # Remove comments from a single file
  annotr clear ./src         # Remove comments from all files in directory`,
	Args: cobra.ExactArgs(1),
	RunE: runClear,
}

func init() {
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	target := args[0]
	info, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("failed to access %s: %w", target, err)
	}

	if info.IsDir() {
		return clearDirectory(target)
	}
	return clearFile(target)
}

func clearFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if !parser.IsSupportedFile(absPath) {
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	fmt.Printf("Clearing comments from %s...\n", filepath.Base(path))

	source, err := fileops.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	cleaned, count := removeComments(source, filepath.Ext(absPath))

	if count > 0 {
		if err := fileops.WriteFile(absPath, cleaned); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	fmt.Printf("✓ Removed %d comment blocks\n", count)
	return nil
}

func clearDirectory(dir string) error {
	files, err := fileops.ScanDirectory(dir)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No supported files found in directory.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	processedCount := 0

	for _, file := range files {
		fmt.Printf("Clear comments from %s? (y/n): ", file.Name)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			if err := clearFile(file.Path); err != nil {
				fmt.Printf("Error processing %s: %v\n", file.Name, err)
			} else {
				processedCount++
			}
		} else {
			fmt.Printf("Skipped %s\n", file.Name)
		}
		fmt.Println()
	}

	fmt.Printf("Done! Cleared comments from %d of %d files.\n", processedCount, len(files))
	return nil
}

func removeComments(source []byte, ext string) ([]byte, int) {
	content := string(source)
	lines := strings.Split(content, "\n")
	var result []string
	count := 0
	inBlockComment := false

	lineCommentPattern := getLineCommentPattern(ext)
	blockStart, blockEnd := getBlockCommentPatterns(ext)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if blockStart != "" && blockEnd != "" {
			if strings.HasPrefix(trimmed, blockStart) && strings.HasSuffix(trimmed, blockEnd) {
				count++
				continue
			}
			if strings.HasPrefix(trimmed, blockStart) {
				inBlockComment = true
				count++
				continue
			}
			if inBlockComment {
				if strings.HasSuffix(trimmed, blockEnd) || trimmed == blockEnd {
					inBlockComment = false
				}
				continue
			}
		}

		if lineCommentPattern != nil && lineCommentPattern.MatchString(trimmed) {
			count++
			continue
		}

		result = append(result, line)
	}

	cleaned := strings.Join(result, "\n")
	cleaned = removeExcessiveBlankLines(cleaned)

	return []byte(cleaned), count
}

func getLineCommentPattern(ext string) *regexp.Regexp {
	switch ext {
	case ".go", ".js", ".ts", ".tsx":
		return regexp.MustCompile(`^//`)
	case ".py":
		return regexp.MustCompile(`^#`)
	default:
		return nil
	}
}

func getBlockCommentPatterns(ext string) (string, string) {
	switch ext {
	case ".go", ".js", ".ts", ".tsx":
		return "/*", "*/"
	case ".py":
		return `"""`, `"""`
	default:
		return "", ""
	}
}

func removeExcessiveBlankLines(content string) string {
	re := regexp.MustCompile(`\n{3,}`)
	return re.ReplaceAllString(content, "\n\n")
}
