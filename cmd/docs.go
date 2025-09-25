package cmd

import (
	"fmt"

	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Open documentation links",
	Long:  `Open relevant documentation in your default browser.`,
}

// docsGeminiCmd opens Gemini CLI documentation
var docsGeminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Open Gemini CLI documentation",
	Long: `Open the Gemini CLI documentation for setup and authentication help.
This includes instructions for API key configuration and usage examples.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("📖 Gemini CLI Documentation", "Setup and Authentication Guide"))
		deps.OpenGeminiDocs()
		return nil
	},
}

// docsBrewCmd opens Homebrew website
var docsBrewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Open Homebrew website",
	Long: `Open the Homebrew website for manual installation instructions.
Use this if you prefer to install Homebrew manually rather than using the automatic installer.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("🍺 Homebrew Website", "Manual Installation Guide"))
		deps.OpenHomebrewSite()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.AddCommand(docsGeminiCmd)
	docsCmd.AddCommand(docsBrewCmd)
}