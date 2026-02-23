## Purpose
Define the Slack mini summary format, derivation rules, and output file conventions.

## Requirements

### Requirement: Slack mini summary is derived from validated main summary
The system SHALL generate a Slack mini summary by parsing the validated main summary text and extracting named sections. The Slack output SHALL contain identical content to the corresponding sections in the main summary with zero modification.

#### Scenario: Main summary contains highlights, action items, and meeting recording
- **WHEN** the validated main summary contains Highlights, Action Items, and Meeting Recording sections
- **THEN** the Slack mini summary includes those sections with byte-identical content (headers and body)

#### Scenario: Main summary contains a Risks section
- **WHEN** the validated main summary contains a Risks section
- **THEN** the Slack mini summary includes the Risks section with identical content

#### Scenario: Main summary does not contain a Risks section
- **WHEN** the validated main summary has no Risks section
- **THEN** the Slack mini summary omits the Risks section without error

### Requirement: Slack mini summary uses fixed section order
The Slack mini summary SHALL present sections in the following fixed order: Title, Highlights, Action Items, Risks, Meeting Recording, Full Meeting Summary. Topic sections from the main summary SHALL be omitted.

#### Scenario: Slack output is assembled from parsed sections
- **WHEN** the section parser extracts sections from the validated main summary
- **THEN** the Slack output contains only: Title, Highlights, Action Items, Risks (if present), Meeting Recording, and Full Meeting Summary — in that exact order

#### Scenario: Topic sections are excluded
- **WHEN** the main summary contains topic sections (italic-only headers)
- **THEN** the Slack mini summary does not include any topic sections

### Requirement: Slack mini summary includes a Full Meeting Summary link
The Slack mini summary SHALL always include a trailing `*FULL MEETING SUMMARY*` section with a placeholder link for the shared drive URL.

#### Scenario: Full Meeting Summary section is appended
- **WHEN** the Slack mini summary is assembled
- **THEN** the final section is `*FULL MEETING SUMMARY*` containing a single bullet with `[Full Summary on Shared Drive](PLACEHOLDER_URL)`

### Requirement: Slack output filename appends -slack suffix
The Slack output filename SHALL be the main summary filename with `-slack` inserted before the `.md` extension.

#### Scenario: Main summary filename includes date and customer
- **WHEN** the main summary filename is `2026-02-23-GISC-cadence-call-summary.md`
- **THEN** the Slack filename is `2026-02-23-GISC-cadence-call-summary-slack.md`

#### Scenario: Main summary filename has no date
- **WHEN** the main summary filename is `GISC-cadence-call-summary.md`
- **THEN** the Slack filename is `GISC-cadence-call-summary-slack.md`

### Requirement: Section parser detects headers by Slack markdown patterns
The section parser SHALL identify section boundaries using the Slack-compatible markdown header patterns defined in the LLM instructions: `*_TITLE_*` for bold-italic title, `_TOPIC_` for italic-only topic headers, and `*HEADER*` for bold-only section headers (Highlights, Action Items, Risks, Meeting Recording).

#### Scenario: Parser encounters a bold-italic title line
- **WHEN** a line matches the pattern `*_..._*`
- **THEN** it is classified as the title section

#### Scenario: Parser encounters a bold-only section header
- **WHEN** a line matches the pattern `*SECTION NAME*` (e.g., `*HIGHLIGHTS*`, `*ACTION ITEMS*`)
- **THEN** it starts a new named section with that header

#### Scenario: Parser encounters an italic-only topic header
- **WHEN** a line matches the pattern `_TOPIC NAME_` (not bold-italic)
- **THEN** it starts a topic section that will be excluded from Slack output
