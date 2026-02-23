## Purpose
Define provider-agnostic AI CLI invocation behavior using `ai.command` with optional `ai.args`.

## Requirements

### Requirement: AI CLI arguments are configurable as ordered tokens
The system SHALL support optional `ai.args` configuration as an ordered list of CLI argument tokens that are passed to the configured AI command without provider-specific interpretation.

#### Scenario: Provider-specific arguments are configured
- **WHEN** settings define `ai.command` and a non-empty `ai.args` list
- **THEN** runtime invocation uses the configured command with the configured arguments in the same order

### Requirement: Empty argument configuration preserves existing no-arg behavior
The system SHALL invoke only the configured AI command when `ai.args` is not configured or is empty.

#### Scenario: Argument list is omitted
- **WHEN** settings define `ai.command` and omit `ai.args`
- **THEN** summary generation invokes the configured command with no additional configured arguments

### Requirement: Legacy inline command arguments remain supported
The system SHALL continue to support whitespace-delimited inline arguments in `ai.command` when `ai.args` is empty.

#### Scenario: Inline command arguments are used without args list
- **WHEN** `ai.command` contains an executable plus inline arguments and `ai.args` is empty
- **THEN** invocation resolves executable and arguments from `ai.command`

### Requirement: Mixed argument sources are rejected
The system MUST reject ambiguous configuration when `ai.command` includes inline arguments and `ai.args` is non-empty.

#### Scenario: Inline arguments and argument list are both configured
- **WHEN** settings define `ai.command` with inline arguments and also define a non-empty `ai.args`
- **THEN** preflight fails with actionable guidance to configure arguments in exactly one location

### Requirement: Preflight validates the resolved executable
The system SHALL validate command availability using the executable resolved from configured invocation values.

#### Scenario: Resolved executable is missing from PATH
- **WHEN** configured invocation resolves to an executable not found in `PATH`
- **THEN** dependency checks fail before summary generation and report guidance that names the missing executable

### Requirement: Runtime execution uses resolved invocation with stdin prompt flow
The system SHALL execute the resolved command invocation while passing prompt content through stdin and capturing stdout/stderr.

#### Scenario: Summary generation runs with configured args
- **WHEN** summary generation executes with configured `ai.command` and `ai.args`
- **THEN** the AI process receives prompt content on stdin while output is captured for summary processing

