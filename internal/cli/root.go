package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "annotr",
	Short: "Fast local code commenting CLI",
	Long: `annotr - Fast Local Code Commenting CLI

Automatically add intelligent, contextual comments to your code files
using local LLM inference via Ollama or cloud providers.

Examples:
  annotr init          # First-time configuration
  annotr file.go       # Add comments to a single file
  annotr ./src         # Process all files in directory`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(updateModelsCmd)
}
