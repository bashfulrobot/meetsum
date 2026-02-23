## Purpose
Define deterministic CI, lint, and local quality baseline expectations for the repository.

## Requirements

### Requirement: CI Go version aligns with module baseline
The CI workflow SHALL use a Go toolchain version compatible with the version declared in `go.mod`.

#### Scenario: CI workflow is run for pull request
- **WHEN** CI initializes Go toolchain
- **THEN** the configured Go version matches the module baseline contract and is not behind declared project requirements

### Requirement: Build artifacts represent real builds
The build workflow SHALL only publish artifacts that are produced by real builds for declared targets.

#### Scenario: Multi-target build workflow runs
- **WHEN** matrix build jobs complete
- **THEN** each published artifact corresponds to an actual built binary for its declared target and no placeholder artifact is published

### Requirement: External workflow actions are version-pinned
The CI and security workflows SHALL use stable, version-pinned action references.

#### Scenario: Workflow action references are evaluated
- **WHEN** repository workflows are reviewed or executed
- **THEN** action references use explicit version tags rather than mutable branches

### Requirement: Lint configuration excludes deprecated linters
The lint baseline SHALL avoid deprecated or removed linter entries.

#### Scenario: Lint job executes
- **WHEN** lint configuration is loaded
- **THEN** lint execution runs without deprecated linter references and produces deterministic rule behavior

### Requirement: Local quality workflow is deterministic
The repository SHALL provide a deterministic local quality command that enforces formatting, vetting, and tests.

#### Scenario: Developer runs local quality command
- **WHEN** the standard local quality command is executed
- **THEN** it runs formatting, vet, and tests in a deterministic order and fails on any check failure

### Requirement: Critical runtime behavior is covered by tests
The repository SHALL include behavior-focused tests for runtime preflight and summary generation safety flows.

#### Scenario: Test suite is executed
- **WHEN** repository tests run
- **THEN** tests include coverage for configured-command preflight outcomes and summary validation-before-write failure paths
