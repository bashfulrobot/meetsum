package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/app"
)

func TestTranscriptContractConsistencyAcrossRuntimeAndValidate(t *testing.T) {
	cfg := newValidationTestConfig(t)
	runtimeService := app.NewService(cfg, nil)

	testCases := []struct {
		name              string
		files             []string
		expectPass        bool
		expectErrorSubset string
	}{
		{
			name:       "exactly one transcript candidate passes",
			files:      []string{"call-notes.TXT"},
			expectPass: true,
		},
		{
			name:              "no transcript candidates fail",
			files:             []string{"notes.md"},
			expectPass:        false,
			expectErrorSubset: "no transcript candidate found",
		},
		{
			name:              "multiple transcript candidates fail with deterministic list",
			files:             []string{"zeta.txt", "alpha.txt"},
			expectPass:        false,
			expectErrorSubset: "alpha.txt, zeta.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			meetingDir := t.TempDir()
			for _, file := range tc.files {
				writeFile(t, filepath.Join(meetingDir, file), "test content")
			}

			_, runtimeErr := runtimeService.Prepare(app.RunRequest{
				UserName:   "Tester",
				MeetingDir: meetingDir,
			})

			results := buildMeetingDirectoryValidationResults(meetingDir, cfg.Files.PovInput)
			validatePass := results[0].Found

			if tc.expectPass {
				if runtimeErr != nil {
					t.Fatalf("expected runtime path to pass, got error: %v", runtimeErr)
				}
				if !validatePass {
					t.Fatalf("expected validate path to pass, got failure detail: %s", results[0].Path)
				}
				return
			}

			if runtimeErr == nil {
				t.Fatalf("expected runtime path to fail")
			}
			if validatePass {
				t.Fatalf("expected validate path to fail")
			}

			runtimeMessage := runtimeErr.Error()
			validateMessage := results[0].Path
			if !strings.Contains(runtimeMessage, tc.expectErrorSubset) {
				t.Fatalf("runtime error %q did not include %q", runtimeMessage, tc.expectErrorSubset)
			}
			if !strings.Contains(validateMessage, tc.expectErrorSubset) {
				t.Fatalf("validate error %q did not include %q", validateMessage, tc.expectErrorSubset)
			}
		})
	}
}

func newValidationTestConfig(t *testing.T) *config.Config {
	t.Helper()

	automationDir := t.TempDir()
	instructionsFile := "instructions.md"
	writeFile(t, filepath.Join(automationDir, instructionsFile), "instructions")

	cfg := &config.Config{}
	cfg.Paths.AutomationDir = automationDir
	cfg.Paths.InstructionsFile = instructionsFile
	cfg.Files.PovInput = "pov-input.md"
	cfg.AI.Command = "echo"

	return cfg
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}
