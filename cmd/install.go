package cmd

import (
	"fmt"

	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure dependencies",
	Long:  `Install and configure required dependencies like gemini-cli and Homebrew.`,
}

// installBrewCmd installs Homebrew
var installBrewCmd = &cobra.Command{
	Use:   "brew",
	Short: "Install Homebrew package manager",
	Long: `Install Homebrew package manager using the official installation script.

⚠️  Security Note: This downloads and executes a shell script from the internet.
For manual installation, run 'meetsum docs brew' to open the Homebrew website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("🍺 Homebrew Installation", "Package Manager for macOS and Linux"))

		if deps.CheckBrewInstalled() {
			fmt.Println(ui.RenderSuccess("✅ Homebrew is already installed!"))
			return nil
		}

		fmt.Println(ui.RenderInfo("💡 Prefer manual installation?"))
		fmt.Println(ui.RenderInfo("Run 'meetsum docs brew' to open installation instructions"))
		fmt.Println()

		return deps.InstallHomebrew()
	},
}

// installGeminiCmd installs gemini-cli
var installGeminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Install gemini-cli",
	Long:  `Install gemini-cli for AI-powered meeting summary generation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("🤖 Gemini CLI Installation", "AI-Powered Text Generation"))

		if deps.CheckGeminiInstalled() {
			fmt.Println(ui.RenderSuccess("✅ gemini-cli is already installed!"))

			// Still offer to open docs for configuration
			fmt.Println()
			fmt.Println(ui.RenderInfo("💡 Need help with configuration?"))
			deps.OpenGeminiDocs()
			return nil
		}

		return deps.InstallGemini()
	},
}

// installAllCmd installs all dependencies
var installAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Install all required dependencies",
	Long:  `Install all required dependencies including Homebrew and gemini-cli.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("🛠️ Dependency Installation", "Setting up meetsum requirements"))

		// Check what's already installed
		brewInstalled := deps.CheckBrewInstalled()
		geminiInstalled := deps.CheckGeminiInstalled()

		if brewInstalled && geminiInstalled {
			fmt.Println(ui.RenderSuccess("✅ All dependencies are already installed!"))
			return nil
		}

		// Install Homebrew if needed
		if !brewInstalled {
			fmt.Println(ui.RenderInfo("📦 Installing Homebrew..."))
			if err := deps.InstallHomebrew(); err != nil {
				return fmt.Errorf("failed to install Homebrew: %w", err)
			}
			fmt.Println()
		} else {
			fmt.Println(ui.RenderSuccess("✅ Homebrew already installed"))
		}

		// Install gemini-cli if needed
		if !geminiInstalled {
			fmt.Println(ui.RenderInfo("🤖 Installing gemini-cli..."))
			if err := deps.InstallGemini(); err != nil {
				return fmt.Errorf("failed to install gemini-cli: %w", err)
			}
		} else {
			fmt.Println(ui.RenderSuccess("✅ gemini-cli already installed"))
		}

		fmt.Println()
		fmt.Println(ui.RenderSuccess("🎉 All dependencies installed successfully!"))
		fmt.Println(ui.RenderInfo("💡 Run 'meetsum' to start generating meeting summaries"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Add subcommands
	installCmd.AddCommand(installBrewCmd)
	installCmd.AddCommand(installGeminiCmd)
	installCmd.AddCommand(installAllCmd)
}