package cmd

import (
	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current configuration settings",
	Long: `Display the current configuration settings in a structured table format.
Shows all configuration values, their defaults, and descriptions.`,
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// Prepare configuration items for display
	configItems := []ui.ConfigItem{
		{
			Category:    "Paths",
			Setting:     "file_browser_root_dir",
			Value:       config.AppConfig.Paths.FileBrowserRootDir,
			Default:     "~/Documents/Company/Customers",
			Description: "Root directory for the file browser",
		},
		{
			Category:    "Paths",
			Setting:     "automation_dir",
			Value:       config.AppConfig.Paths.AutomationDir,
			Default:     "~/Documents/Company/automation/summaries",
			Description: "Directory containing LLM instructions",
		},
		{
			Category:    "Paths",
			Setting:     "instructions_file",
			Value:       config.AppConfig.Paths.InstructionsFile,
			Default:     "Meeting-summary-llm-instructions.md",
			Description: "Name of the AI instructions file",
		},
		{
			Category:    "Files",
			Setting:     "transcript",
			Value:       config.AppConfig.Files.Transcript,
			Default:     "transcript.txt",
			Description: "Required transcript filename in meeting directories",
		},
		{
			Category:    "Files",
			Setting:     "pov_input",
			Value:       config.AppConfig.Files.PovInput,
			Default:     "pov-input.md",
			Description: "Optional context file for additional meeting details",
		},
		{
			Category:    "AI",
			Setting:     "command",
			Value:       config.AppConfig.AI.Command,
			Default:     "gemini",
			Description: "AI CLI command for text generation",
		},
		{
			Category:    "Features",
			Setting:     "trace_mode",
			Value:       boolToString(config.AppConfig.Features.TraceMode),
			Default:     "false",
			Description: "Enable detailed output and disable loading indicators",
		},
		{
			Category:    "Features",
			Setting:     "file_browser",
			Value:       boolToString(config.AppConfig.Features.FileBrowser),
			Default:     "true",
			Description: "Enable interactive file picker for directory selection",
		},
		{
			Category:    "User",
			Setting:     "name",
			Value:       config.AppConfig.User.Name,
			Default:     "(not set)",
			Description: "Default name for first-person perspective summaries",
		},
	}

	// Display the configuration table
	return ui.ShowConfigurationTable(configItems)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func init() {
	rootCmd.AddCommand(configCmd)
}