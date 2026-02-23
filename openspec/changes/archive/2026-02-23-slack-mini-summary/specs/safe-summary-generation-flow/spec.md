## MODIFIED Requirements

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
