package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/deps"
	"github.com/bashfulrobot/meetsum/internal/summary"
	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/charmbracelet/glamour"
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
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
		Level:           log.InfoLevel,
	})

	if traceMode || config.AppConfig.Features.TraceMode {
		logger.SetLevel(log.DebugLevel)
	}
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

		// Offer to select different directory
		var tryAgain bool
		err := huh.NewConfirm().
			Title("Would you like to select a different directory?").
			Value(&tryAgain).
			Run()
		if err != nil {
			return err
		}

		if tryAgain {
			fmt.Println(ui.RenderInfo("🗂️  Opening file picker..."))
			meetingDir, err = ui.SelectDirectory(config.AppConfig.Paths.CustomersDir)
			if err != nil {
				return err
			}
			processor.SetMeetingDir(meetingDir)
			if err := processor.ValidateRequiredFiles(); err != nil {
				fmt.Println(ui.RenderError("Still no transcript.txt found. Exiting."))
				return err
			}
		} else {
			return fmt.Errorf("required files not found")
		}
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

	// Confirm before proceeding
	fmt.Println()
	var proceed bool
	err = huh.NewConfirm().
		Title("Generate meeting summary?").
		Value(&proceed).
		Run()
	if err != nil {
		return err
	}

	if !proceed {
		fmt.Println(ui.RenderInfo("👋 Operation cancelled"))
		return nil
	}

	// Show progress and generate summary
	fmt.Println()
	fmt.Println(ui.RenderInfo("🚀 Preparing to generate summary..."))

	// Final confirmation with model info
	fmt.Println()
	fmt.Println(ui.RenderInfoBox(
		"🤖 Model: Gemini Pro",
		fmt.Sprintf("📍 Working Directory: %s", meetingDir),
		"⚡ Ready to generate summary",
	))

	var startProcessing bool
	err = huh.NewConfirm().
		Title("Start Gemini Pro processing? (This may take several minutes)").
		Value(&startProcessing).
		Run()
	if err != nil {
		return err
	}

	if !startProcessing {
		fmt.Println(ui.RenderInfo("👋 Operation cancelled"))
		return nil
	}

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

	// Use progress bar for the multi-step operation unless in trace mode
	var summary string
	if config.AppConfig.Features.TraceMode {
		fmt.Println(ui.RenderInfo("🧠 Gemini Pro is processing your meeting transcript..."))
		summary, err = processor.GenerateSummary()
	} else {
		steps := []string{
			"Reading transcript file",
			"Loading AI instructions",
			"Processing with Gemini Pro",
			"Generating formatted summary",
		}

		result, err := ui.RunWithProgress(steps, func(updateProgress func(step int, message string)) (interface{}, error) {
			updateProgress(0, "Reading transcript and context files")
			// Small delay to show progress step

			updateProgress(1, "Loading LLM instructions")
			// Another step

			updateProgress(2, "Sending to Gemini Pro for processing")
			summary, err := processor.GenerateSummary()
			if err != nil {
				return nil, err
			}

			updateProgress(3, "Formatting and finalizing summary")
			return summary, nil
		})

		if err != nil {
			return err
		}
		summary = result.(string)
	}

	if err != nil {
		fmt.Println(ui.RenderError(fmt.Sprintf("Failed to generate summary: %v", err)))
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

	// Show success message
	fmt.Println(ui.RenderInfoBox(
		fmt.Sprintf("📄 Summary file: %s", filepath.Base(outputPath)),
		fmt.Sprintf("📍 Location: %s", meetingDir),
	))

	// Offer to view the summary
	if config.AppConfig.Features.MarkdownPreview {
		var viewSummary bool
		err = huh.NewConfirm().
			Title("Would you like to preview the generated summary?").
			Value(&viewSummary).
			Run()
		if err == nil && viewSummary {
			fmt.Println(ui.AccentStyle.Render("📖 Summary Preview:"))
			fmt.Println()

			// Use glamour for markdown rendering
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(80),
			)
			if err == nil {
				rendered, err := renderer.Render(summary)
				if err == nil {
					// Limit preview to reasonable length
					lines := strings.Split(rendered, "\n")
					if len(lines) > 50 {
						fmt.Print(strings.Join(lines[:50], "\n"))
						fmt.Println(ui.WarningStyle.Render("... (truncated - full summary in file)"))
					} else {
						fmt.Print(rendered)
					}
				} else {
					// Fallback to plain text
					lines := strings.Split(summary, "\n")
					if len(lines) > 50 {
						fmt.Println(strings.Join(lines[:50], "\n"))
						fmt.Println(ui.WarningStyle.Render("... (truncated - full summary in file)"))
					} else {
						fmt.Println(summary)
					}
				}
			}
		}
	}

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
		// Validate provided directory
		if _, err := os.Stat(meetingDir); os.IsNotExist(err) {
			return "", fmt.Errorf("directory '%s' does not exist", meetingDir)
		}
		return meetingDir, nil
	}

	fmt.Println(ui.RenderInfo("📁 Enter the meeting directory path:"))
	if config.AppConfig.Features.FileBrowser {
		fmt.Println(ui.SecondaryStyle.Render("   (or press Enter for guided selection)"))
	}

	var inputPath string
	err := huh.NewInput().
		Title("Meeting Directory Path").
		Description("Directory should contain transcript.txt").
		Placeholder("~/Documents/Customers/[Customer]/[date] or press Enter").
		Value(&inputPath).
		Run()

	if err != nil {
		return "", err
	}

	inputPath = strings.TrimSpace(inputPath)

	// If no path entered and file browser is enabled, use file picker
	if inputPath == "" && config.AppConfig.Features.FileBrowser {
		fmt.Println(ui.RenderInfo("🗂️  Opening file picker..."))
		return ui.SelectDirectory(config.AppConfig.Paths.CustomersDir)
	}

	if inputPath == "" {
		return "", fmt.Errorf("no directory selected")
	}

	// Validate the directory
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory '%s' does not exist", inputPath)
	}

	return inputPath, nil
}

