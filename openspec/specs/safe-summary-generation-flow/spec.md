## Purpose
Define safe summary generation behavior with explicit execution context and guarded output persistence.

## Requirements

### Requirement: Summary generation avoids process-wide working-directory mutation
The summary generation flow SHALL execute provider commands with explicit directory context and SHALL NOT mutate process-global working directory state.

#### Scenario: Summary generation is executed for a meeting directory
- **WHEN** a summary generation run is started
- **THEN** provider execution uses explicit meeting-directory context without calling process-wide directory mutation for the full process

### Requirement: Final summary write happens only after output validation
The system SHALL validate generated summary content before writing the final summary file. After the main summary is persisted, the system SHALL also generate and persist the Slack mini summary. A failure to write the Slack file SHALL NOT fail the run — it SHALL be surfaced as a non-fatal warning.

#### Scenario: Generated summary output is empty or invalid
- **WHEN** cleaned model output fails summary validation
- **THEN** the final summary file is not written, the Slack file is not written, and the run fails with an actionable validation error

#### Scenario: Main summary is valid and saved successfully
- **WHEN** the main summary passes validation and is written to disk
- **THEN** the Slack mini summary is generated from the validated content and written to the meeting directory

#### Scenario: Slack file write fails after main summary succeeds
- **WHEN** the main summary is written successfully but the Slack file write fails
- **THEN** the run succeeds with the main summary output and a non-fatal warning about the Slack file failure is surfaced in the result

### Requirement: Invalid output preserves diagnostics
The system SHALL preserve raw provider output for diagnostics when summary validation fails.

#### Scenario: Validation fails after provider execution
- **WHEN** summary output fails validation
- **THEN** raw provider output is persisted to a diagnostic artifact location and surfaced in error guidance

### Requirement: Transcript discovery behavior is preserved
The modernization flow SHALL preserve transcript-source requirements defined in `discover-transcript-from-txt`.

#### Scenario: Meeting directory transcript selection is evaluated
- **WHEN** runtime generation validates transcript source
- **THEN** transcript selection follows the same `0/1/many .txt` confidence contract established by the transcript discovery change
