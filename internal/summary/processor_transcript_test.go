package summary

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bashfulrobot/meetsum/config"
)

func TestFindTranscriptFile(t *testing.T) {
	// Load config with defaults
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Run("finds standard transcript.txt", func(t *testing.T) {
		testDir := t.TempDir()
		transcriptPath := filepath.Join(testDir, "transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)

		found, err := processor.FindTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if found != transcriptPath {
			t.Errorf("Expected %s, got %s", transcriptPath, found)
		}
	})

	t.Run("finds dated transcript when standard missing", func(t *testing.T) {
		testDir := t.TempDir()
		datedPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if err := os.WriteFile(datedPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)

		found, err := processor.FindTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if found != datedPath {
			t.Errorf("Expected %s, got %s", datedPath, found)
		}
	})

	t.Run("prefers standard transcript.txt over dated", func(t *testing.T) {
		testDir := t.TempDir()
		standardPath := filepath.Join(testDir, "transcript.txt")
		datedPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if err := os.WriteFile(standardPath, []byte("standard"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := os.WriteFile(datedPath, []byte("dated"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)

		found, err := processor.FindTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if found != standardPath {
			t.Errorf("Expected standard path %s, got %s", standardPath, found)
		}
	})

	t.Run("ignores invalid date format", func(t *testing.T) {
		testDir := t.TempDir()
		// Invalid: wrong format
		invalidPath := filepath.Join(testDir, "02-04-2026-transcript.txt")
		if err := os.WriteFile(invalidPath, []byte("invalid"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)

		_, err := processor.FindTranscriptFile()
		if err == nil {
			t.Error("Expected error for missing valid transcript")
		}
	})

	t.Run("returns error when no transcript found", func(t *testing.T) {
		testDir := t.TempDir()

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)

		_, err := processor.FindTranscriptFile()
		if err == nil {
			t.Error("Expected error when no transcript found")
		}
	})
}

func TestRenameTranscriptFile(t *testing.T) {
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Run("renames transcript.txt to dated format", func(t *testing.T) {
		// Create structure: tempdir/2026-02-04/transcript.txt
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}
		transcriptPath := filepath.Join(testDir, "transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if newName != "2026-02-04-transcript.txt" {
			t.Errorf("Expected '2026-02-04-transcript.txt', got '%s'", newName)
		}

		// Verify file was renamed
		newPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Error("Renamed file does not exist")
		}
		if _, err := os.Stat(transcriptPath); !os.IsNotExist(err) {
			t.Error("Original file still exists")
		}
	})

	t.Run("skips rename for already dated transcript", func(t *testing.T) {
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}
		datedPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if err := os.WriteFile(datedPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if newName != "" {
			t.Errorf("Expected empty string (no rename), got '%s'", newName)
		}
	})

	t.Run("skips rename when no date in folder path", func(t *testing.T) {
		testDir := t.TempDir() // No date in path
		transcriptPath := filepath.Join(testDir, "transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		newName, err := processor.RenameTranscriptFile()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if newName != "" {
			t.Errorf("Expected empty string (no rename), got '%s'", newName)
		}

		// Verify file was NOT renamed
		if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
			t.Error("Original file should still exist")
		}
	})

	t.Run("errors when destination exists", func(t *testing.T) {
		baseDir := t.TempDir()
		testDir := filepath.Join(baseDir, "2026-02-04")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}
		// Create both files
		transcriptPath := filepath.Join(testDir, "transcript.txt")
		datedPath := filepath.Join(testDir, "2026-02-04-transcript.txt")
		if err := os.WriteFile(transcriptPath, []byte("original"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := os.WriteFile(datedPath, []byte("existing"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		processor := NewProcessor(config.AppConfig, nil)
		processor.SetMeetingDir(testDir)
		if err := processor.ValidateRequiredFiles(); err != nil {
			t.Fatalf("ValidateRequiredFiles failed: %v", err)
		}

		_, err := processor.RenameTranscriptFile()
		if err == nil {
			t.Error("Expected error when destination exists")
		}
	})
}
