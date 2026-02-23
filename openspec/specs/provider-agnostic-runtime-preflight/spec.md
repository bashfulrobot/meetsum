## Purpose
Define provider-agnostic runtime dependency preflight behavior driven by configured AI command.

## Requirements

### Requirement: Runtime preflight uses configured AI command
The system SHALL perform runtime preflight against the configured `ai.command` value instead of hard-coded provider command names.

#### Scenario: Configured command is available
- **WHEN** `ai.command` is set to a command that exists in `PATH`
- **THEN** runtime preflight succeeds without requiring a provider-specific fallback check

### Requirement: Runtime preflight is provider-agnostic
The system SHALL NOT block runtime summary generation because of missing Gemini-specific dependencies when `ai.command` is configured to a different provider command.

#### Scenario: Non-Gemini provider command is configured
- **WHEN** `ai.command` is set to a non-Gemini executable and that executable is available
- **THEN** runtime preflight allows summary generation without requiring Homebrew or Gemini installation

### Requirement: Missing configured command returns actionable guidance
The system SHALL return an actionable runtime preflight error when the configured `ai.command` executable is not available.

#### Scenario: Configured command is missing
- **WHEN** `ai.command` points to an executable not found in `PATH`
- **THEN** runtime preflight fails with an error that includes the configured command name and clear next steps

### Requirement: Dependency reporting reflects configured command
The system SHALL report dependency readiness based on configured runtime provider command status.

#### Scenario: Dependency check command is executed
- **WHEN** the user runs dependency checks
- **THEN** the output includes configured provider command status and does not incorrectly report readiness for a missing configured command
