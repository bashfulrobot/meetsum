package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [directory]",
	Short: "Validate meeting directory structure and files",
	Long: `Validate a meeting directory to ensure it contains all required files
and check for optional context files. If no directory is specified,
validates the configuration files and directory structure.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	var targetDir string

	if len(args) > 0 {
		targetDir = args[0]
		return validateMeetingDirectory(targetDir)
	} else {
		return validateConfiguration()
	}
}

func validateMeetingDirectory(meetingDir string) error {
	// Expand home directory if needed
	if strings.HasPrefix(meetingDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			meetingDir = filepath.Join(homeDir, meetingDir[2:])
		}
	}

	// Check if directory exists
	if _, err := os.Stat(meetingDir); os.IsNotExist(err) {
		return err
	}

	// Prepare file validation results
	results := []ui.FileValidationResult{
		{
			File:        config.AppConfig.Files.Transcript,
			Required:    true,
			Description: "Meeting transcript file - required for processing",
		},
		{
			File:        config.AppConfig.Files.PovInput,
			Required:    false,
			Description: "Point of view input file - optional context",
		},
	}

	// Check each file
	for i := range results {
		filePath := filepath.Join(meetingDir, results[i].File)
		if _, err := os.Stat(filePath); err == nil {
			results[i].Found = true
			results[i].Path = filePath
		} else {
			results[i].Found = false
			results[i].Path = ""
		}
	}

	// Also check for other common files
	commonFiles := []string{
		"notes.md",
		"agenda.md",
		"recording.mp4",
		"recording.m4a",
		"attendees.txt",
	}

	for _, file := range commonFiles {
		filePath := filepath.Join(meetingDir, file)
		if _, err := os.Stat(filePath); err == nil {
			results = append(results, ui.FileValidationResult{
				File:        file,
				Required:    false,
				Found:       true,
				Path:        filePath,
				Description: "Additional meeting file found",
			})
		}
	}

	return ui.ShowFileValidationTable(results)
}

func validateConfiguration() error {
	// Validate configuration files and directories
	results := []ui.FileValidationResult{
		{
			File:        "Instructions File",
			Required:    true,
			Path:        config.AppConfig.GetInstructionsPath(),
			Description: "AI instructions for summary generation",
		},
		{
			File:        "Customers Directory",
			Required:    true,
			Path:        config.AppConfig.Paths.CustomersDir,
			Description: "Base directory for customer meeting folders",
		},
		{
			File:        "Automation Directory",
			Required:    true,
			Path:        config.AppConfig.Paths.AutomationDir,
			Description: "Directory containing automation files",
		},
	}

	// Check each path/file
	for i := range results {
		if _, err := os.Stat(results[i].Path); err == nil {
			results[i].Found = true
		} else {
			results[i].Found = false
		}
	}

	return ui.ShowFileValidationTable(results)
}

func init() {
	rootCmd.AddCommand(validateCmd)
}