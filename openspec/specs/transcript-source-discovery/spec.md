## Purpose
Define transcript source discovery behavior from meeting-directory text files with deterministic confidence handling.

## Requirements

### Requirement: Transcript candidates are discovered from text files
The system SHALL discover transcript candidates from files in the selected meeting directory whose filename extension is `.txt` (case-insensitive).

#### Scenario: Exactly one transcript candidate exists
- **WHEN** a meeting directory contains exactly one file with a `.txt` extension
- **THEN** the system selects that file as the transcript source

### Requirement: Missing transcript candidates fail fast
The system SHALL fail summary processing and validation when a meeting directory contains no transcript candidates.

#### Scenario: No text files are present
- **WHEN** a meeting directory contains zero files with a `.txt` extension
- **THEN** the system reports that no transcript candidate was found and does not continue

### Requirement: Ambiguous transcript candidates fail fast
The system SHALL fail summary processing and validation when a meeting directory contains multiple transcript candidates.

#### Scenario: Multiple text files are present
- **WHEN** a meeting directory contains more than one file with a `.txt` extension
- **THEN** the system reports an ambiguity error and does not continue

### Requirement: Ambiguity feedback is deterministic
The system SHALL include the discovered transcript candidate filenames in deterministic order when reporting ambiguity.

#### Scenario: Candidate list is displayed in error output
- **WHEN** an ambiguity error is returned for multiple transcript candidates
- **THEN** the error output includes all candidate filenames in a stable order

### Requirement: Transcript discovery behavior is consistent across command paths
The system SHALL apply the same transcript-discovery and confidence rules in summary generation, directory selection validation, and the `validate` command.

#### Scenario: Directory acceptance and runtime processing use the same transcript contract
- **WHEN** a user provides or selects a meeting directory
- **THEN** runtime generation, picker validation, and `validate` produce equivalent pass/fail outcomes for transcript discovery

