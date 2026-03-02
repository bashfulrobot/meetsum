package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ai"
	"github.com/bashfulrobot/meetsum/internal/app"
	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	traceMode  bool
	askName    bool
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

// SetVersion sets the version information for the application.
// Also sets rootCmd.Version so cobra's built-in --version flag works.
func SetVersion(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
	rootCmd.Version = v
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/meetsum/settings.yaml)")
	rootCmd.Flags().BoolVar(&traceMode, "trace", false, "Run without spinners to see all output")
	rootCmd.Flags().BoolVar(&askName, "ask-name", false, "Prompt for name even if default is configured")
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

	runtimeService := app.NewService(config.AppConfig, logger)
	aiCommand, err := runtimeService.Preflight()
	if err != nil {
		fmt.Println(ui.RenderHeader("🤖 Meeting Summary Generator", "Runtime Preflight"))
		fmt.Println(ui.RenderError(err.Error()))
		fmt.Println()
		fmt.Println(ui.RenderInfo("💡 Next steps:"))
		for _, line := range preflightGuidance(aiCommand) {
			fmt.Println(ui.RenderInfo("  • " + line))
		}
		return err
	}

	// Display header
	fmt.Println(ui.RenderHeader("🤖 Meeting Summary Generator", fmt.Sprintf("Powered by %s", aiCommand)))

	// Get user's name - use config default unless --ask-name is set
	var userName string
	defaultName := strings.TrimSpace(config.AppConfig.User.Name)
	if defaultName != "" && !askName {
		userName = defaultName
		fmt.Println(ui.RenderInfo(fmt.Sprintf("👤 Using configured name: %s", userName)))
	} else {
		userName, err = getUserName()
		if err != nil {
			return err
		}
	}

	// Get meeting directory
	if len(args) > 0 {
		meetingDir = args[0]
	}

	meetingDir, err = getMeetingDirectory()
	if err != nil {
		return err
	}

	session, err := runtimeService.Prepare(app.RunRequest{
		UserName:   userName,
		MeetingDir: meetingDir,
	})
	if err != nil {
		fmt.Println(ui.RenderError(err.Error()))
		return err
	}
	preparation := session.Preparation()

	// Show summary of found files
	fmt.Println(ui.RenderInfoBox(
		fmt.Sprintf("📁 Meeting Directory: %s", filepath.Base(preparation.MeetingDir)),
		fmt.Sprintf("📄 Transcript: ✅ %s", preparation.TranscriptFile),
		"📋 Instructions: ✅ Found",
	))

	// Check for optional files
	if len(preparation.OptionalFiles) > 0 {
		fmt.Println(ui.RenderSuccess("🎯 Context files found:"))
		for _, file := range preparation.OptionalFiles {
			fmt.Println(ui.FileListStyle.Render("  " + file))
		}
	} else {
		fmt.Println(ui.RenderWarning("No context files found (pov-input.md)"))
	}

	// Show processing info
	_, resolvedArgs, resolveErr := ai.ResolveConfiguredInvocation(config.AppConfig.AI.Command, config.AppConfig.AI.Args)
	if resolveErr != nil {
		return resolveErr
	}

	fmt.Println()
	fmt.Printf("🤖 Runtime Command: %s\n", aiCommand)
	fmt.Printf("🧩 Runtime Args: %d configured token(s)\n", len(resolvedArgs))
	fmt.Printf("📍 Working Directory: %s\n", preparation.MeetingDir)
	fmt.Println("⚡ Starting summary generation...")
	fmt.Println()

	var runResult app.RunResult
	if config.AppConfig.Features.TraceMode {
		fmt.Printf("🧠 %s is processing your meeting transcript...\n", aiCommand)
		runResult, err = session.Run()
	} else {
		result, err := ui.RunWithSpinner(
			fmt.Sprintf("🧠 %s is processing your meeting transcript...", aiCommand),
			func() (any, error) {
				return session.Run()
			},
		)
		if err != nil {
			return err
		}

		typedResult, ok := result.(app.RunResult)
		if !ok {
			return fmt.Errorf("unexpected runtime result type: %T", result)
		}
		runResult = typedResult
	}

	if err != nil {
		fmt.Println(ui.RenderError(fmt.Sprintf("Failed to generate summary: %v", err)))
		if config.AppConfig.Logging.Output == "file" || config.AppConfig.Logging.Output == "both" {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("💡 Check the log file for detailed error output: %s", config.AppConfig.GetLogFilePath())))
		}
		return err
	}

	// Show success message
	infoLines := []string{
		fmt.Sprintf("📄 Summary file: %s", filepath.Base(runResult.OutputPath)),
		fmt.Sprintf("📍 Location: %s", preparation.MeetingDir),
	}
	if runResult.SlackOutputPath != "" {
		infoLines = append(infoLines, fmt.Sprintf("📋 Slack summary: %s", filepath.Base(runResult.SlackOutputPath)))
	}
	if runResult.RenamedTranscript != "" {
		infoLines = append(infoLines, fmt.Sprintf("📝 Transcript renamed to: %s", runResult.RenamedTranscript))
	}
	fmt.Println(ui.RenderInfoBox(infoLines...))
	if runResult.RenameWarning != "" {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("Could not rename transcript: %s", runResult.RenameWarning)))
	}
	if runResult.SlackWarning != "" {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("Could not save Slack summary: %s", runResult.SlackWarning)))
	}

	fmt.Println()
	fmt.Println(ui.RenderSuccess("🎉 All done! Your meeting summary is ready."))

	return nil
}

func preflightGuidance(aiCommand string) []string {
	command := strings.TrimSpace(aiCommand)
	if command == "" {
		resolved, _, err := ai.ResolveConfiguredInvocation(config.AppConfig.AI.Command, config.AppConfig.AI.Args)
		if err == nil {
			command = resolved
		}
	}
	if command == "" {
		command = "<configured-ai-command>"
	}

	lines := []string{
		fmt.Sprintf("install %q and ensure it is available in PATH", command),
		"update ai.command/ai.args in your settings.yaml to a valid invocation for this machine",
		"run 'meetsum check' to re-validate runtime dependencies",
	}

	if command == "gemini" {
		if deps.CheckBrewInstalled() {
			lines = append(lines, "run 'meetsum install gemini' for guided Gemini CLI installation")
		} else {
			lines = append(lines, "run 'meetsum install all' to install Homebrew and Gemini CLI")
		}
	}

	return lines
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
		Description("Directory must contain exactly one .txt transcript file").
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
