package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ai"
	"github.com/bashfulrobot/meetsum/internal/summary"
	"github.com/charmbracelet/log"
)

// RunRequest captures runtime inputs collected by CLI handlers.
type RunRequest struct {
	UserName   string
	MeetingDir string
}

// RunPreparation captures validated runtime context for CLI rendering.
type RunPreparation struct {
	MeetingDir     string
	TranscriptFile string
	OptionalFiles  []string
}

// RunResult captures output from a runtime summary execution.
type RunResult struct {
	Summary           string
	OutputPath        string
	SlackOutputPath   string
	RenamedTranscript string
	RenameWarning     string
	SlackWarning      string
}

// Service orchestrates summary runtime behavior independently of CLI rendering.
type Service struct {
	cfg    *config.Config
	logger *log.Logger
}

// Session encapsulates a prepared runtime execution.
type Session struct {
	processor   *summary.Processor
	preparation RunPreparation
}

// NewService creates a runtime service.
func NewService(cfg *config.Config, logger *log.Logger) *Service {
	return &Service{cfg: cfg, logger: logger}
}

// Preflight checks configured runtime dependencies.
func (s *Service) Preflight() (string, error) {
	command, err := ai.CheckConfiguredCommandAvailable(s.cfg.AI.Command, s.cfg.AI.Args)
	if err != nil {
		return command, err
	}
	return command, nil
}

// Prepare validates runtime inputs and required files.
func (s *Service) Prepare(request RunRequest) (*Session, error) {
	userName := strings.TrimSpace(request.UserName)
	if userName == "" {
		return nil, fmt.Errorf("user name is required")
	}

	meetingDir := strings.TrimSpace(request.MeetingDir)
	if meetingDir == "" {
		return nil, fmt.Errorf("meeting directory is required")
	}

	processor := summary.NewProcessor(s.cfg, s.logger)
	processor.SetUserName(userName)
	processor.SetMeetingDir(meetingDir)

	if err := processor.ValidateRequiredFiles(); err != nil {
		return nil, err
	}

	preparation := RunPreparation{
		MeetingDir:     meetingDir,
		TranscriptFile: filepath.Base(processor.TranscriptPath()),
		OptionalFiles:  processor.GetOptionalFiles(),
	}

	return &Session{
		processor:   processor,
		preparation: preparation,
	}, nil
}

// Preparation returns session context suitable for CLI display.
func (s *Session) Preparation() RunPreparation {
	return s.preparation
}

// Run executes summary generation and persistence.
func (s *Session) Run() (RunResult, error) {
	output, err := s.processor.GenerateSummaryOutput()
	if err != nil {
		return RunResult{}, err
	}

	if err := s.processor.ValidateSummaryContent(output.Cleaned); err != nil {
		diagnosticPath, saveErr := s.processor.SaveRawOutputDiagnostics(output.Raw)
		if saveErr != nil {
			return RunResult{}, fmt.Errorf("%w; also failed to save diagnostics: %w", err, saveErr)
		}

		return RunResult{}, fmt.Errorf(
			"%w; raw provider output saved to %s",
			err,
			diagnosticPath,
		)
	}

	outputPath, err := s.processor.SaveSummary(output.Cleaned)
	if err != nil {
		return RunResult{}, err
	}

	// Generate and save Slack mini summary (non-fatal)
	slackOutputPath := ""
	slackWarning := ""
	sections := summary.ParseSections(output.Cleaned)
	slackContent := summary.BuildSlackSummary(sections)
	slackPath, slackErr := s.processor.SaveSlackSummary(slackContent)
	if slackErr != nil {
		slackWarning = slackErr.Error()
	} else {
		slackOutputPath = slackPath
	}

	renamedTranscript, err := s.processor.RenameTranscriptFile()
	renameWarning := ""
	if err != nil {
		renameWarning = err.Error()
	}

	return RunResult{
		Summary:           output.Cleaned,
		OutputPath:        outputPath,
		SlackOutputPath:   slackOutputPath,
		RenamedTranscript: renamedTranscript,
		RenameWarning:     renameWarning,
		SlackWarning:      slackWarning,
	}, nil
}
