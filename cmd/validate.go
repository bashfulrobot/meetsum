package cmd

import (
	"os"
	"path/filepath"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/summary"
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
	if len(args) > 0 {
		return validateMeetingDirectory(args[0])
	}
	return validateConfiguration()
}

func validateMeetingDirectory(meetingDir string) error {
	meetingDir = expandPath(meetingDir)

	// Check if directory exists
	if _, err := os.Stat(meetingDir); os.IsNotExist(err) {
		return err
	}

	results := buildMeetingDirectoryValidationResults(meetingDir, config.AppConfig.Files.PovInput)
	return ui.ShowFileValidationTable(results)
}

func buildMeetingDirectoryValidationResults(meetingDir, povInputFile string) []ui.FileValidationResult {
	// Prepare file validation results
	results := []ui.FileValidationResult{
		{
			File:        "*.txt transcript",
			Required:    true,
			Description: "Exactly one .txt transcript file is required for processing",
		},
		{
			File:        povInputFile,
			Required:    false,
			Description: "Point of view input file - optional context",
		},
	}

	// Check transcript discovery contract first.
	transcriptPath, err := summary.FindSingleTranscriptCandidate(meetingDir)
	if err == nil {
		results[0].Found = true
		results[0].Path = transcriptPath
		results[0].Description = "Transcript candidate selected"
	} else {
		results[0].Found = false
		results[0].Path = err.Error()
		results[0].Description = "Expected exactly one .txt transcript candidate"
	}

	// Check optional POV input file.
	povPath := filepath.Join(meetingDir, results[1].File)
	if _, err := os.Stat(povPath); err == nil {
		results[1].Found = true
		results[1].Path = povPath
	} else {
		results[1].Found = false
		results[1].Path = ""
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

	return results
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
			File:        "File Browser Root",
			Required:    true,
			Path:        config.AppConfig.Paths.FileBrowserRootDir,
			Description: "Root directory for the file browser",
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
