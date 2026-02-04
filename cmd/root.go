package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/summary"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	traceMode  bool
	meetingDir string
	cfgFile    string
	logger     *log.Logger

	// Version information
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meetsum [meeting_directory]",
	Short: "Meeting Summary Generator CLI Tool",
	Long: `A CLI tool that generates structured meeting summaries using AI.
Reads meeting transcripts and generates formatted summaries with context.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMeetSum,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// SetVersion sets the version information for the application
func SetVersion(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/meetsum/settings.yaml)")
	rootCmd.Flags().BoolVar(&traceMode, "trace", false, "Run without spinners to see all output")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		fmt.Printf("Using config file: %s\n", cfgFile)
	}

	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		// Note: cobra.OnInitialize functions can't return errors,
		// but we can't exit here either. This is a limitation we need to address.
		// For now, we'll continue with defaults but log the error.
	}
}

func initLogger() {
	cfg := config.AppConfig

	// Parse log level
	level := log.InfoLevel
	switch strings.ToLower(cfg.Logging.Level) {
	case "debug":
		level = log.DebugLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	}

	// Override with trace mode
	if traceMode || cfg.Features.TraceMode {
		level = log.DebugLevel
	}

	// Determine output writer(s)
	var writer io.Writer
	switch cfg.Logging.Output {
	case "file":
		writer = openLogFile(cfg.GetLogFilePath())
	case "both":
		writer = io.MultiWriter(os.Stderr, openLogFile(cfg.GetLogFilePath()))
	default: // "screen"
		writer = os.Stderr
	}

	logger = log.NewWithOptions(writer, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		Level:           level,
	})
}

func openLogFile(path string) io.Writer {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create log directory: %v\n", err)
		return os.Stderr
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not open log file: %v\n", err)
		return os.Stderr
	}
	return file
}

func runMeetSum(cmd *cobra.Command, args []string) error {
	// Override trace mode from config if flag is set
	if cmd.Flags().Changed("trace") {
		config.AppConfig.Features.TraceMode = traceMode
	}

	// Check dependencies - only gemini-cli is required
	geminiInstalled := deps.CheckGeminiInstalled()

	if !geminiInstalled {
		fmt.Println(ui.RenderHeader("🤖 Meeting Summary Generator", "Dependency Check"))

		brewInstalled := deps.CheckBrewInstalled()
		if !brewInstalled {
			// Need to install Homebrew first, then gemini-cli
			fmt.Println(ui.RenderError("gemini-cli is required but not installed"))
			fmt.Println(ui.RenderWarning("Homebrew is needed to install gemini-cli"))
			fmt.Println()
			fmt.Println(ui.RenderInfo("💡 Installation options:"))
			fmt.Println(ui.RenderInfo("  • Run 'meetsum install all' to install all dependencies"))
			fmt.Println(ui.RenderInfo("  • Run 'meetsum install brew' then 'meetsum install gemini'"))
			fmt.Println(ui.RenderInfo("  • Run 'meetsum check' to verify installation"))
			fmt.Println()

			var installNow bool
			err := huh.NewConfirm().
				Title("Would you like to install all dependencies now?").
				Description("This will install Homebrew first, then gemini-cli").
				Value(&installNow).
				Run()
			if err != nil {
				return err
			}

			if installNow {
				// Install Homebrew first
				if err := deps.InstallHomebrew(); err != nil {
					return fmt.Errorf("failed to install Homebrew: %w", err)
				}
				fmt.Println()

				// Then install gemini-cli
				if err := deps.InstallGemini(); err != nil {
					return fmt.Errorf("failed to install gemini-cli: %w", err)
				}

				fmt.Println()
				fmt.Println(ui.RenderSuccess("✅ All dependencies installed! You can now run meetsum again."))
				return nil
			} else {
				return fmt.Errorf("gemini-cli is required to generate summaries")
			}
		} else {
			// Homebrew available, just need gemini-cli
			fmt.Println(ui.RenderError("gemini-cli is required but not installed"))
			fmt.Println()
			fmt.Println(ui.RenderInfo("💡 Installation options:"))
			fmt.Println(ui.RenderInfo("  • Run 'meetsum install gemini' to install via Homebrew"))
			fmt.Println(ui.RenderInfo("  • Run 'meetsum check' to verify installation"))
			fmt.Println()

			var installNow bool
			err := huh.NewConfirm().
				Title("Would you like to install gemini-cli now?").
				Value(&installNow).
				Run()
			if err != nil {
				return err
			}

			if installNow {
				if err := deps.InstallGemini(); err != nil {
					return fmt.Errorf("failed to install gemini-cli: %w", err)
				}
				fmt.Println()
				fmt.Println(ui.RenderSuccess("✅ Installation complete! You can now run meetsum again."))
				return nil
			} else {
				return fmt.Errorf("gemini-cli is required to generate summaries")
			}
		}
	}

	// Display header
	fmt.Println(ui.RenderHeader("🤖 Meeting Summary Generator", "Powered by Gemini Pro"))

	// Get user's name
	userName, err := getUserName()
	if err != nil {
		return err
	}

	// Get meeting directory
	if len(args) > 0 {
		meetingDir = args[0]
	}

	meetingDir, err = getMeetingDirectory()
	if err != nil {
		return err
	}

	// Initialize processor
	processor := summary.NewProcessor(config.AppConfig, logger)
	processor.SetUserName(userName)
	processor.SetMeetingDir(meetingDir)

	// Validate required files
	if err := processor.ValidateRequiredFiles(); err != nil {
		fmt.Println(ui.RenderError(err.Error()))
		return err
	}

	// Show summary of found files
	fmt.Println(ui.RenderInfoBox(
		fmt.Sprintf("📁 Meeting Directory: %s", filepath.Base(meetingDir)),
		"📄 Transcript: ✅ Found",
		"📋 Instructions: ✅ Found",
	))

	// Check for optional files
	optionalFiles := processor.GetOptionalFiles()
	if len(optionalFiles) > 0 {
		fmt.Println(ui.RenderSuccess("🎯 Context files found:"))
		for _, file := range optionalFiles {
			fmt.Println(ui.FileListStyle.Render("  " + file))
		}
	} else {
		fmt.Println(ui.RenderWarning("No context files found (pov-input.md)"))
	}

	// Show processing info
	fmt.Println()
	fmt.Println("🤖 Model: Gemini Pro")
	fmt.Printf("📍 Working Directory: %s\n", meetingDir)
	fmt.Println("⚡ Starting summary generation...")

	// Generate summary with spinner
	fmt.Println()

	// Change to meeting directory for proper path context
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(meetingDir); err != nil {
		return err
	}

	// Use spinner for processing unless in trace mode
	var summary string
	if config.AppConfig.Features.TraceMode {
		fmt.Println(ui.RenderInfo("🧠 Gemini Pro is processing your meeting transcript..."))
		summary, err = processor.GenerateSummary()
	} else {
		result, err := ui.RunWithSpinner("🧠 Gemini Pro is processing your meeting transcript...", func() (interface{}, error) {
			return processor.GenerateSummary()
		})

		if err != nil {
			return err
		}
		summary = result.(string)
	}

	if err != nil {
		fmt.Println(ui.RenderError(fmt.Sprintf("Failed to generate summary: %v", err)))
		if config.AppConfig.Logging.Output == "file" || config.AppConfig.Logging.Output == "both" {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("💡 Check the log file for detailed error output: %s", config.AppConfig.GetLogFilePath())))
		}
		return err
	}

	// Save summary
	outputPath, err := processor.SaveSummary(summary)
	if err != nil {
		fmt.Println(ui.RenderError(fmt.Sprintf("Failed to save summary: %v", err)))
		return err
	}

	// Check if summary was generated successfully
	if summary == "" {
		fmt.Println(ui.RenderError("No output generated. Check gemini-cli installation and authentication."))
		return fmt.Errorf("empty summary generated")
	}

	// Rename transcript file to dated format if applicable
	renamedTo, renameErr := processor.RenameTranscriptFile()
	if renameErr != nil {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("Could not rename transcript: %v", renameErr)))
	}

	// Show success message
	infoLines := []string{
		fmt.Sprintf("📄 Summary file: %s", filepath.Base(outputPath)),
		fmt.Sprintf("📍 Location: %s", meetingDir),
	}
	if renamedTo != "" {
		infoLines = append(infoLines, fmt.Sprintf("📝 Transcript renamed to: %s", renamedTo))
	}
	fmt.Println(ui.RenderInfoBox(infoLines...))


	fmt.Println()
	fmt.Println(ui.RenderSuccess("🎉 All done! Your meeting summary is ready."))

	return nil
}

func getUserName() (string, error) {
	fmt.Println(ui.RenderInfo("👤 Enter your name (for first-person perspective):"))

	var userName string
	err := huh.NewInput().
		Title("Your Name").
		Placeholder("Enter your name").
		Value(&userName).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("name is required")
			}
			return nil
		}).
		Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(userName), nil
}

func getMeetingDirectory() (string, error) {
	if meetingDir != "" {
		// Expand ~ to home directory
		expandedDir := expandPath(meetingDir)

		// Validate provided directory
		if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
			return "", fmt.Errorf("directory '%s' does not exist", meetingDir)
		}

		// Convert to absolute path for proper metadata extraction
		absPath, err := filepath.Abs(expandedDir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		return absPath, nil
	}

	// If file browser is enabled, go directly to file picker
	if config.AppConfig.Features.FileBrowser {
		startDir := expandPath(config.AppConfig.Paths.FileBrowserRootDir)
		selectedPath, err := ui.SelectDirectory(startDir)
		if err != nil {
			return "", err
		}
		// File picker returns paths that should already be absolute, but ensure it
		absPath, err := filepath.Abs(selectedPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		return absPath, nil
	}

	// File browser disabled - prompt for manual path entry
	fmt.Println(ui.RenderInfo("📁 Enter the meeting directory path:"))

	var inputPath string
	err := huh.NewInput().
		Title("Meeting Directory Path").
		Description("Directory should contain transcript.txt").
		Placeholder("~/Documents/Customers/[Customer]/[date]").
		Value(&inputPath).
		Run()

	if err != nil {
		return "", err
	}

	inputPath = strings.TrimSpace(inputPath)

	if inputPath == "" {
		return "", fmt.Errorf("no directory path provided")
	}

	// Expand ~ to home directory
	expandedPath := expandPath(inputPath)

	// Validate the directory
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory '%s' does not exist", inputPath)
	}

	// Convert to absolute path for proper metadata extraction
	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return absPath, nil
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

