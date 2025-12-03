package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	tmpPath := filepath.Join(dir, ".annotr-tmp-"+filename)

	if err := os.WriteFile(tmpPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func InsertComment(source []byte, lineNum uint32, comment string, language string) []byte {
	lines := strings.Split(string(source), "\n")
	
	if int(lineNum) > len(lines) {
		return source
	}

	indent := getIndent(lines[lineNum])
	commentLines := strings.Split(comment, "\n")
	for i, cl := range commentLines {
		commentLines[i] = indent + cl
	}
	formattedComment := strings.Join(commentLines, "\n")

	newLines := make([]string, 0, len(lines)+len(commentLines))
	newLines = append(newLines, lines[:lineNum]...)
	newLines = append(newLines, formattedComment)
	newLines = append(newLines, lines[lineNum:]...)

	return []byte(strings.Join(newLines, "\n"))
}

func getIndent(line string) string {
	var indent strings.Builder
	for _, ch := range line {
		if ch == ' ' || ch == '\t' {
			indent.WriteRune(ch)
		} else {
			break
		}
	}
	return indent.String()
}
