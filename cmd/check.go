package cmd

import (
	"fmt"
	"os/exec"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ai"
	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check system dependencies and configuration",
	Long: `Check if all required dependencies are installed and properly configured.
This includes verifying that the configured AI command is available and functional.`,
	RunE: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderHeader("🔍 Dependency Check", "Verifying meetsum requirements"))

	allGood := true
	configured := config.AppConfig.AI.Command
	resolved, resolveErr := ai.ResolveCommand(configured)

	fmt.Printf("🤖 configured ai.command (%q): ", configured)
	if resolveErr != nil {
		fmt.Println(ui.RenderError("❌ Invalid configuration"))
		fmt.Println(ui.RenderInfo(fmt.Sprintf("   %v", resolveErr)))
		allGood = false
	} else {
		if _, err := ai.CheckCommandAvailable(configured); err != nil {
			fmt.Println(ui.RenderError("❌ Not available"))
			fmt.Println(ui.RenderInfo(fmt.Sprintf("   Install %q or update ai.command in settings.yaml", resolved)))
			allGood = false
		} else {
			fmt.Println(ui.RenderSuccess("✅ Available"))
		}
	}

	if resolved == "gemini" {
		fmt.Print("🔧 gemini setup: ")
		if err := deps.ValidateGeminiSetup(); err != nil {
			fmt.Println(ui.RenderWarning("⚠️  May need configuration"))
			fmt.Println(ui.RenderInfo("   Run 'meetsum install gemini' for setup help"))
		} else {
			fmt.Println(ui.RenderSuccess("✅ Functional"))
		}

		fmt.Print("🍺 Homebrew: ")
		if deps.CheckBrewInstalled() {
			fmt.Println(ui.RenderSuccess("✅ Installed"))
		} else {
			fmt.Println(ui.RenderWarning("⚠️  Not installed (needed for guided gemini install path)"))
		}
	}

	fmt.Print("📋 git: ")
	if _, err := exec.LookPath("git"); err == nil {
		fmt.Println(ui.RenderSuccess("✅ Available"))
	} else {
		fmt.Println(ui.RenderWarning("⚠️  Not found (optional)"))
	}

	fmt.Println()

	if allGood {
		fmt.Println(ui.RenderSuccess("🎉 Runtime dependencies are ready!"))
		fmt.Println(ui.RenderInfo("💡 You can now run 'meetsum' to generate meeting summaries"))
	} else {
		fmt.Println(ui.RenderError("❌ Required runtime dependency checks failed"))
		fmt.Println(ui.RenderInfo("💡 Install the configured AI command or update ai.command"))
	}

	fmt.Println()
	fmt.Println(ui.RenderInfo("📊 Use 'meetsum config' to view detailed configuration"))

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
