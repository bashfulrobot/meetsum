package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ai"
)

func TestServicePreflightMissingCommand(t *testing.T) {
	cfg := newTestConfig(t, "missing-ai-command")
	service := NewService(cfg, nil)

	command, err := service.Preflight()
	if err == nil {
		t.Fatalf("expected preflight error for missing command")
	}

	if command != "missing-ai-command" {
		t.Fatalf("expected command name to be propagated, got %q", command)
	}
	if !errors.Is(err, ai.ErrMissingCommand) {
		t.Fatalf("expected ErrMissingCommand, got: %v", err)
	}
}

func TestServiceRunSuccessAndRename(t *testing.T) {
	commandDir := t.TempDir()
	writeExecutable(t, commandDir, "fake-ai-success", `#!/usr/bin/env bash
cat >/dev/null
cat <<'OUT'
*_SUMMARY_*
- Completed action items
OUT
`)
	t.Setenv("PATH", commandDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfg := newTestConfig(t, "fake-ai-success")
	service := NewService(cfg, nil)

	command, err := service.Preflight()
	if err != nil {
		t.Fatalf("preflight failed unexpectedly: %v", err)
	}
	if command != "fake-ai-success" {
		t.Fatalf("expected fake-ai-success command, got %q", command)
	}

	meetingDir := createMeetingDir(t, "2026-02-04", "transcript.txt", "transcript content")
	session, err := service.Prepare(RunRequest{UserName: "Tester", MeetingDir: meetingDir})
	if err != nil {
		t.Fatalf("prepare failed: %v", err)
	}

	preparation := session.Preparation()
	if preparation.TranscriptFile != "transcript.txt" {
		t.Fatalf("expected transcript.txt preparation, got %q", preparation.TranscriptFile)
	}

	result, err := session.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if result.OutputPath == "" {
		t.Fatalf("expected output path")
	}
	if result.RenamedTranscript != "2026-02-04-transcript.txt" {
		t.Fatalf("expected transcript rename, got %q", result.RenamedTranscript)
	}

	content, err := os.ReadFile(result.OutputPath)
	if err != nil {
		t.Fatalf("failed to read output summary: %v", err)
	}
	if !strings.Contains(string(content), "*_SUMMARY_*") {
		t.Fatalf("expected cleaned summary content, got: %s", string(content))
	}

	renamedPath := filepath.Join(meetingDir, result.RenamedTranscript)
	if _, err := os.Stat(renamedPath); err != nil {
		t.Fatalf("expected renamed transcript file: %v", err)
	}
}

func TestServiceRunCommandFailure(t *testing.T) {
	commandDir := t.TempDir()
	writeExecutable(t, commandDir, "fake-ai-error", `#!/usr/bin/env bash
cat >/dev/null
echo "simulated failure" >&2
exit 1
`)
	t.Setenv("PATH", commandDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfg := newTestConfig(t, "fake-ai-error")
	service := NewService(cfg, nil)
	meetingDir := createMeetingDir(t, "2026-02-04", "transcript.txt", "transcript content")

	session, err := service.Prepare(RunRequest{UserName: "Tester", MeetingDir: meetingDir})
	if err != nil {
		t.Fatalf("prepare failed: %v", err)
	}

	_, err = session.Run()
	if err == nil {
		t.Fatalf("expected command execution failure")
	}
	if !strings.Contains(err.Error(), "failed to generate summary") {
		t.Fatalf("expected generate-summary failure, got: %v", err)
	}

	summaryMatches, err := filepath.Glob(filepath.Join(meetingDir, "*-cadence-call-summary.md"))
	if err != nil {
		t.Fatalf("failed to glob summary files: %v", err)
	}
	if len(summaryMatches) != 0 {
		t.Fatalf("expected no persisted summary file, found %v", summaryMatches)
	}
}

func TestServiceValidationFailurePersistsRawOutput(t *testing.T) {
	commandDir := t.TempDir()
	writeExecutable(t, commandDir, "fake-ai-invalid", `#!/usr/bin/env bash
cat >/dev/null
echo "Loaded cached credentials."
`)
	t.Setenv("PATH", commandDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfg := newTestConfig(t, "fake-ai-invalid")
	service := NewService(cfg, nil)
	meetingDir := createMeetingDir(t, "2026-02-04", "transcript.txt", "transcript content")

	session, err := service.Prepare(RunRequest{UserName: "Tester", MeetingDir: meetingDir})
	if err != nil {
		t.Fatalf("prepare failed: %v", err)
	}

	_, err = session.Run()
	if err == nil {
		t.Fatalf("expected validation failure for empty cleaned output")
	}
	if !strings.Contains(err.Error(), "raw provider output saved to") {
		t.Fatalf("expected diagnostic guidance in error, got: %v", err)
	}

	diagnosticPath := filepath.Join(meetingDir, "summary-raw-output.txt")
	diagnosticContent, err := os.ReadFile(diagnosticPath)
	if err != nil {
		t.Fatalf("expected diagnostic output file: %v", err)
	}
	if !strings.Contains(string(diagnosticContent), "Loaded cached credentials.") {
		t.Fatalf("expected raw diagnostic content, got: %s", string(diagnosticContent))
	}

	summaryMatches, err := filepath.Glob(filepath.Join(meetingDir, "*-cadence-call-summary.md"))
	if err != nil {
		t.Fatalf("failed to glob summary files: %v", err)
	}
	if len(summaryMatches) != 0 {
		t.Fatalf("expected no summary file on validation failure, found %v", summaryMatches)
	}
}

func newTestConfig(t *testing.T, command string) *config.Config {
	t.Helper()

	automationDir := t.TempDir()
	instructionsFile := "instructions.md"
	instructionsPath := filepath.Join(automationDir, instructionsFile)
	if err := os.WriteFile(instructionsPath, []byte("Meeting instructions"), 0644); err != nil {
		t.Fatalf("failed to write instructions file: %v", err)
	}

	cfg := &config.Config{}
	cfg.Paths.AutomationDir = automationDir
	cfg.Paths.InstructionsFile = instructionsFile
	cfg.Files.PovInput = "pov-input.md"
	cfg.AI.Command = command

	return cfg
}

func createMeetingDir(t *testing.T, meetingDate, transcriptName, transcriptContent string) string {
	t.Helper()

	baseDir := t.TempDir()
	meetingDir := filepath.Join(baseDir, "Customers", "Acme", meetingDate)
	if err := os.MkdirAll(meetingDir, 0755); err != nil {
		t.Fatalf("failed to create meeting dir: %v", err)
	}

	transcriptPath := filepath.Join(meetingDir, transcriptName)
	if err := os.WriteFile(transcriptPath, []byte(transcriptContent), 0644); err != nil {
		t.Fatalf("failed to write transcript file: %v", err)
	}

	return meetingDir
}

func writeExecutable(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write test executable %s: %v", name, err)
	}
}
