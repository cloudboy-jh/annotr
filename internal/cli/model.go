package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudboy-jh/annotr/internal/config"
	"github.com/cloudboy-jh/annotr/internal/ui"
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Change the default model",
	Long: `Change the default model used for generating comments.

This will show available models for your configured provider.`,
	RunE: runModel,
}

func init() {
	rootCmd.AddCommand(modelCmd)
}

func runModel(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		fmt.Println("No configuration found. Run 'annotr init' first.")
		return nil
	}

	model := ui.NewModelSelectModel(cfg)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running model selector: %w", err)
	}

	if m, ok := finalModel.(ui.ModelSelectModel); ok {
		if m.Error() != nil {
			os.Exit(1)
		}
	}

	return nil
}
