package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnhorton/annotr/internal/config"
	"github.com/johnhorton/annotr/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize annotr configuration",
	Long: `Initialize annotr configuration.

This will guide you through:
  - Detecting Ollama or selecting a cloud provider
  - Selecting a model
  - Choosing a comment style

Configuration is saved to ~/.annotr/config.json`,
	RunE: runInit,
}

// // Allows user to create or update configuration by running the initialization program.
// // If configuration already exists, prompts user to delete it before proceeding.
func runInit(cmd *cobra.Command, args []string) error {
	if config.Exists() {
		fmt.Println("Configuration already exists at ~/.annotr/config.json")
		fmt.Println("Delete it first if you want to reconfigure.")
		return nil
	}

	model := ui.NewInitModel()
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running init: %w", err)
	}

	if m, ok := finalModel.(ui.InitModel); ok {
		if m.Error() != nil {
			os.Exit(1)
		}
	}

	return nil
}
