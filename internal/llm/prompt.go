package llm

import (
	"fmt"
	"strings"
)

type CommentTarget struct {
	Language     string
	Filename     string
	Code         string
	Context      string
	CommentStyle string
}

func BuildCommentPrompt(target CommentTarget) []Message {
	systemPrompt := `You are a code documentation expert. Generate concise, accurate comments for code blocks.
Rules:
- Be brief but informative
- Focus on the "why" not the "what"
- Use the specified comment style
- Return ONLY the comment text, no code
- Do not include comment delimiters (like // or /* */)
- Maximum 1-2 sentences for simple functions
- Maximum 3-4 sentences for complex logic`

	userPrompt := fmt.Sprintf(`Language: %s
File: %s
Comment Style: %s

Context:
%s

Target Code:
%s

Generate a comment for the target code:`, 
		target.Language,
		target.Filename,
		target.CommentStyle,
		target.Context,
		target.Code,
	)

	return []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
}

func FormatComment(comment, language, style string) string {
	comment = strings.TrimSpace(comment)
	
	switch style {
	case "line":
		return formatLineComment(comment, language)
	case "block":
		return formatBlockComment(comment, language)
	case "docstring":
		return formatDocstring(comment, language)
	default:
		return formatLineComment(comment, language)
	}
}

func formatLineComment(comment, language string) string {
	prefix := getLineCommentPrefix(language)
	lines := strings.Split(comment, "\n")
	var result []string
	for _, line := range lines {
		result = append(result, prefix+" "+strings.TrimSpace(line))
	}
	return strings.Join(result, "\n")
}

func formatBlockComment(comment, language string) string {
	start, end := getBlockCommentDelimiters(language)
	return fmt.Sprintf("%s %s %s", start, comment, end)
}

func formatDocstring(comment, language string) string {
	switch language {
	case "python":
		return fmt.Sprintf(`"""%s"""`, comment)
	case "go":
		lines := strings.Split(comment, "\n")
		var result []string
		for _, line := range lines {
			result = append(result, "// "+strings.TrimSpace(line))
		}
		return strings.Join(result, "\n")
	default:
		return formatBlockComment(comment, language)
	}
}

func getLineCommentPrefix(language string) string {
	switch language {
	case "python", "ruby", "shell", "bash", "yaml":
		return "#"
	default:
		return "//"
	}
}

func getBlockCommentDelimiters(language string) (string, string) {
	switch language {
	case "python":
		return `"""`, `"""`
	case "html", "xml":
		return "<!--", "-->"
	default:
		return "/*", "*/"
	}
}
