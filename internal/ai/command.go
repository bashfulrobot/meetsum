package ai

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrMissingCommand = errors.New("configured AI command is not available")

// MissingCommandError captures configured-command preflight failures.
type MissingCommandError struct {
	Command string
}

func (e MissingCommandError) Error() string {
	return fmt.Sprintf("configured AI command %q was not found in PATH", e.Command)
}

// ResolveCommand extracts the executable name from configured ai.command.
func ResolveCommand(configured string) (string, error) {
	parts := strings.Fields(strings.TrimSpace(configured))
	if len(parts) == 0 {
		return "", fmt.Errorf("ai.command is empty; configure it in settings.yaml")
	}
	return parts[0], nil
}

// ResolveCommandArgs parses configured ai.command into executable and arguments.
func ResolveCommandArgs(configured string) (string, []string, error) {
	parts := strings.Fields(strings.TrimSpace(configured))
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("ai.command is empty; configure it in settings.yaml")
	}
	return parts[0], parts[1:], nil
}

// CheckCommandAvailable verifies the configured command exists in PATH.
func CheckCommandAvailable(configured string) (string, error) {
	command, err := ResolveCommand(configured)
	if err != nil {
		return "", err
	}

	if _, err := exec.LookPath(command); err != nil {
		return command, errors.Join(MissingCommandError{Command: command}, ErrMissingCommand)
	}

	return command, nil
}
