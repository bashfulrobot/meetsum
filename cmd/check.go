package cmd

import (
	"fmt"
	"os/exec"

	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check system dependencies and configuration",
	Long: `Check if all required dependencies are installed and properly configured.
This includes verifying that gemini-cli is available and functional.`,
	RunE: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderHeader("🔍 Dependency Check", "Verifying meetsum requirements"))

	allGood := true

	// Check Homebrew
	fmt.Print("🍺 Homebrew: ")
	if deps.CheckBrewInstalled() {
		fmt.Println(ui.RenderSuccess("✅ Installed"))
	} else {
		fmt.Println(ui.RenderError("❌ Not installed"))
		allGood = false
	}

	// Check gemini-cli
	fmt.Print("🤖 gemini-cli: ")
	if deps.CheckGeminiInstalled() {
		fmt.Println(ui.RenderSuccess("✅ Installed"))

		// Test if gemini is functional
		fmt.Print("🔧 gemini configuration: ")
		if err := deps.ValidateGeminiSetup(); err != nil {
			fmt.Println(ui.RenderWarning("⚠️  May need configuration"))
			fmt.Println(ui.RenderInfo("   Run 'meetsum install gemini' for setup help"))
		} else {
			fmt.Println(ui.RenderSuccess("✅ Functional"))
		}
	} else {
		fmt.Println(ui.RenderError("❌ Not installed"))
		allGood = false
	}

	// Check other useful tools
	fmt.Print("📋 git: ")
	if _, err := exec.LookPath("git"); err == nil {
		fmt.Println(ui.RenderSuccess("✅ Available"))
	} else {
		fmt.Println(ui.RenderWarning("⚠️  Not found (optional)"))
	}

	fmt.Println()

	if allGood {
		fmt.Println(ui.RenderSuccess("🎉 All dependencies are ready!"))
		fmt.Println(ui.RenderInfo("💡 You can now run 'meetsum' to generate meeting summaries"))
	} else {
		fmt.Println(ui.RenderError("❌ Some dependencies are missing"))
		fmt.Println(ui.RenderInfo("💡 Run 'meetsum install all' to install missing dependencies"))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}