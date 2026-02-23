package summary

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bashfulrobot/meetsum/config"
)

func TestFindTranscriptFile(t *testing.T) {
	t.Run("finds exactly one txt transcript candidate", func(t *testing.T) {
		testDir := t.TempDir()
		expectedPath := filepath.Join(testDir, "meeting-notes.txt")
		if err := os.WriteFile(expectedPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		found, err := processor.FindTranscriptFile()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if found != expectedPath {
			t.Fatalf("expected %s, got %s", expectedPath, found)
		}
	})

	t.Run("treats TXT extension case-insensitively", func(t *testing.T) {
		testDir := t.TempDir()
		expectedPath := filepath.Join(testDir, "Transcript.TXT")
		if err := os.WriteFile(expectedPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		found, err := processor.FindTranscriptFile()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if found != expectedPath {
			t.Fatalf("expected %s, got %s", expectedPath, found)
		}
	})

	t.Run("returns error when no txt transcript candidates exist", func(t *testing.T) {
		testDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(testDir, "notes.md"), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create non-transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		_, err := processor.FindTranscriptFile()
		if err == nil {
			t.Fatalf("expected error when no .txt transcript candidates exist")
		}

		if !strings.Contains(err.Error(), "no transcript candidate found") {
			t.Fatalf("expected no-candidate error, got: %v", err)
		}
	})

	t.Run("returns deterministic candidate list when transcript source is ambiguous", func(t *testing.T) {
		testDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(testDir, "zeta.txt"), []byte("zeta"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(testDir, "alpha.txt"), []byte("alpha"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		_, err := processor.FindTranscriptFile()
		if err == nil {
			t.Fatalf("expected ambiguity error")
		}

		errText := err.Error()
		if !strings.Contains(errText, "alpha.txt, zeta.txt") {
			t.Fatalf("expected sorted candidate list in error, got: %s", errText)
		}
	})
}

func TestRenameTranscriptFile(t *testing.T) {
	t.Run("renames selected transcript to dated format", func(t *testing.T) {
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("failed to create test dir: %v", err)
		}

		transcriptPath := filepath.Join(testDir, "transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if newName != "2026-02-04-transcript.txt" {
			t.Fatalf("expected 2026-02-04-transcript.txt, got %s", newName)
		}

		newPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Fatalf("renamed file does not exist")
		}
		if _, err := os.Stat(transcriptPath); !os.IsNotExist(err) {
			t.Fatalf("original file still exists")
		}
	})

	t.Run("skips rename when transcript is already dated", func(t *testing.T) {
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("failed to create test dir: %v", err)
		}

		datedPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if err := os.WriteFile(datedPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if newName != "" {
			t.Fatalf("expected no rename, got %s", newName)
		}
	})

	t.Run("skips rename when no date exists in meeting path", func(t *testing.T) {
		testDir := t.TempDir()
		transcriptPath := filepath.Join(testDir, "transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if newName != "" {
			t.Fatalf("expected no rename, got %s", newName)
		}

		if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
			t.Fatalf("original transcript file should still exist")
		}
	})

	t.Run("errors when destination transcript already exists", func(t *testing.T) {
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("failed to create test dir: %v", err)
		}

		if err := os.WriteFile(filepath.Join(testDir, "transcript.txt"), []byte("original"), 0644); err != nil {
			t.Fatalf("failed to create transcript file: %v", err)
		}

		processor := newTestProcessor(t, testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		if err := os.WriteFile(filepath.Join(testDir, "2026-02-04-transcript.txt"), []byte("existing"), 0644); err != nil {
			t.Fatalf("failed to create dated transcript file: %v", err)
		}

		_, err := processor.RenameTranscriptFile()
		if err == nil {
			t.Fatalf("expected error when destination exists")
		}
	})
}

func newTestProcessor(t *testing.T, meetingDir string) *Processor {
	t.Helper()

	automationDir := t.TempDir()
	instructionsFile := "instructions.md"
	instructionsPath := filepath.Join(automationDir, instructionsFile)
	if err := os.WriteFile(instructionsPath, []byte("test instructions"), 0644); err != nil {
		t.Fatalf("failed to create instructions file: %v", err)
	}

	cfg := &config.Config{}
	cfg.Paths.AutomationDir = automationDir
	cfg.Paths.InstructionsFile = instructionsFile
	cfg.Files.PovInput = "pov-input.md"
	cfg.AI.Command = "echo"

	processor := NewProcessor(cfg, nil)
	processor.SetMeetingDir(meetingDir)
	processor.SetUserName("Test User")
	return processor
}
