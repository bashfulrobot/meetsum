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

// ResolveConfiguredInvocation resolves executable and args from ai.command + ai.args.
// Compatibility behavior:
// - If ai.args is non-empty, ai.command must contain only the executable token.
// - If ai.args is empty, ai.command may include inline args and is tokenized by whitespace.
func ResolveConfiguredInvocation(configuredCommand string, configuredArgs []string) (string, []string, error) {
	parts := strings.Fields(strings.TrimSpace(configuredCommand))
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("ai.command is empty; configure it in settings.yaml")
	}

	if len(configuredArgs) > 0 {
		if len(parts) > 1 {
			return "", nil, fmt.Errorf(
				"ai.command contains inline arguments while ai.args is also set; keep executable in ai.command and move all flags to ai.args",
			)
		}
		return parts[0], append([]string(nil), configuredArgs...), nil
	}

	return parts[0], parts[1:], nil
}

// ResolveCommand extracts the executable name from configured ai.command.
func ResolveCommand(configured string) (string, error) {
	command, _, err := ResolveConfiguredInvocation(configured, nil)
	if err != nil {
		return "", err
	}
	return command, nil
}

// ResolveCommandArgs parses configured ai.command into executable and arguments.
func ResolveCommandArgs(configured string) (string, []string, error) {
	return ResolveConfiguredInvocation(configured, nil)
}

// CheckConfiguredCommandAvailable verifies the resolved executable exists in PATH.
func CheckConfiguredCommandAvailable(configuredCommand string, configuredArgs []string) (string, error) {
	command, _, err := ResolveConfiguredInvocation(configuredCommand, configuredArgs)
	if err != nil {
		return "", err
	}

	if _, err := exec.LookPath(command); err != nil {
		return command, errors.Join(MissingCommandError{Command: command}, ErrMissingCommand)
	}

	return command, nil
}

// CheckCommandAvailable verifies the configured command exists in PATH.
func CheckCommandAvailable(configured string) (string, error) {
	return CheckConfiguredCommandAvailable(configured, nil)
}
