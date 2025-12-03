package parser

import (
	"strings"
)

func BuildContext(source []byte, block CodeBlock, contextLines int) string {
	lines := strings.Split(string(source), "\n")
	
	startLine := int(block.StartLine)
	if startLine > contextLines {
		startLine = startLine - contextLines
	} else {
		startLine = 0
	}

	endLine := int(block.EndLine) + contextLines
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}

	var contextParts []string
	
	if startLine < int(block.StartLine) {
		beforeLines := lines[startLine:block.StartLine]
		contextParts = append(contextParts, "// Before:\n"+strings.Join(beforeLines, "\n"))
	}

	if int(block.EndLine) < endLine {
		afterLines := lines[block.EndLine+1 : endLine+1]
		if len(afterLines) > 0 {
			contextParts = append(contextParts, "// After:\n"+strings.Join(afterLines, "\n"))
		}
	}

	return strings.Join(contextParts, "\n\n")
}

func ExtractImports(source []byte, language string) string {
	lines := strings.Split(string(source), "\n")
	var imports []string

	switch language {
	case "go":
		inImportBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "import (") {
				inImportBlock = true
				imports = append(imports, line)
			} else if inImportBlock {
				imports = append(imports, line)
				if trimmed == ")" {
					inImportBlock = false
				}
			} else if strings.HasPrefix(trimmed, "import ") {
				imports = append(imports, line)
			}
		}
	case "python":
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "from ") {
				imports = append(imports, line)
			}
		}
	case "javascript", "typescript":
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "const ") && strings.Contains(trimmed, "require(") {
				imports = append(imports, line)
			}
		}
	}

	return strings.Join(imports, "\n")
}
