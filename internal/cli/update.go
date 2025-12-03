package cli

import (
	"fmt"

	"github.com/cloudboy-jh/annotr/internal/config"
	"github.com/spf13/cobra"
)

var updateModelsCmd = &cobra.Command{
	Use:   "update-models",
	Short: "Update the models manifest",
	Long: `Update the models manifest with the latest available models.

This refreshes the list of available models for each provider.`,
	RunE: runUpdateModels,
}

func runUpdateModels(cmd *cobra.Command, args []string) error {
	manifest := config.DefaultModelsManifest()

	if err := manifest.Save(); err != nil {
		return fmt.Errorf("failed to save models manifest: %w", err)
	}

	fmt.Println("✓ Models manifest updated")
	fmt.Println("  Saved to: ~/.annotr/models.json")

	return nil
}
