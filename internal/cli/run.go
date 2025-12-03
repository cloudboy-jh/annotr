package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnhorton/annotr/internal/config"
	"github.com/johnhorton/annotr/internal/fileops"
	"github.com/johnhorton/annotr/internal/llm"
	"github.com/johnhorton/annotr/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Args = cobra.MaximumNArgs(1)
	rootCmd.RunE = runAnnotate
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		fmt.Println("No configuration found. Run 'annotr init' first.")
		return nil
	}

	target := args[0]
	info, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("failed to access %s: %w", target, err)
	}

	if info.IsDir() {
		return processDirectory(cfg, target)
	}
	return processFile(cfg, target)
}

func processFile(cfg *config.Config, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if !parser.IsSupportedFile(absPath) {
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	fmt.Printf("Processing %s...\n", filepath.Base(path))

	source, err := fileops.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	p, err := parser.NewParser(absPath)
	if err != nil {
		return err
	}

	blocks, err := p.Parse(source)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	if len(blocks) == 0 {
		fmt.Println("No commentable code blocks found.")
		return nil
	}

	apiKey := cfg.APIKeys[cfg.DefaultProvider]
	client := llm.NewClient(cfg.DefaultProvider, apiKey, cfg.DefaultModel)

	commentCount := 0
	modifiedSource := source

	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]

		if hasExistingComment(modifiedSource, block.StartLine) {
			continue
		}

		ctx := parser.BuildContext(source, block, 5)
		target := llm.CommentTarget{
			Language:     p.Language(),
			Filename:     filepath.Base(absPath),
			Code:         block.Code,
			Context:      ctx,
			CommentStyle: cfg.CommentStyle,
		}

		messages := llm.BuildCommentPrompt(target)
		resp, err := client.Complete(context.Background(), &llm.CompletionRequest{
			Messages:  messages,
			MaxTokens: 256,
		})
		if err != nil {
			fmt.Printf("Warning: failed to generate comment for %s: %v\n", block.Name, err)
			continue
		}

		comment := llm.FormatComment(resp.Content, p.Language(), cfg.CommentStyle)
		modifiedSource = fileops.InsertComment(modifiedSource, block.StartLine, comment, p.Language())
		commentCount++
	}

	if commentCount > 0 {
		if err := fileops.WriteFile(absPath, modifiedSource); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	fmt.Printf("✓ Added %d comments\n\n", commentCount)
	fmt.Println("Enjoy your comments! ;)")

	return nil
}

func processDirectory(cfg *config.Config, dir string) error {
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
	skippedCount := 0

	for _, file := range files {
		fmt.Printf("Process %s? (y/n): ", file.Name)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			if err := processFile(cfg, file.Path); err != nil {
				fmt.Printf("Error processing %s: %v\n", file.Name, err)
			} else {
				processedCount++
			}
		} else {
			fmt.Printf("Skipped %s\n", file.Name)
			skippedCount++
		}
		fmt.Println()
	}

	fmt.Printf("Done! Commented %d of %d files.\n", processedCount, len(files))
	fmt.Println("Enjoy your comments! ;)")

	return nil
}

func hasExistingComment(source []byte, lineNum uint32) bool {
	lines := strings.Split(string(source), "\n")
	if lineNum == 0 {
		return false
	}
	if int(lineNum-1) >= len(lines) {
		return false
	}

	prevLine := strings.TrimSpace(lines[lineNum-1])
	return strings.HasPrefix(prevLine, "//") ||
		strings.HasPrefix(prevLine, "#") ||
		strings.HasPrefix(prevLine, "/*") ||
		strings.HasPrefix(prevLine, "*/") ||
		strings.HasPrefix(prevLine, "*") ||
		strings.HasPrefix(prevLine, `"""`) ||
		strings.HasPrefix(prevLine, "'''")
}
