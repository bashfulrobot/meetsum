## Why

The current workflow assumes a fixed transcript filename, which is brittle for teams that manage many customers and many date folders with varied naming conventions. This also creates behavioral drift because runtime processing, validation, and directory selection do not enforce the same transcript contract.

## What Changes

- Replace fixed transcript filename assumptions with transcript discovery based on `*.txt` files in the selected meeting directory.
- Define deterministic selection behavior: exactly one `*.txt` file uses that file, zero `*.txt` files fails with an actionable error, and more than one `*.txt` file fails with an ambiguity error that lists candidate files.
- Align transcript validation rules across `meetsum` runtime flow, file picker selection checks, and `meetsum validate`.
- **BREAKING**: directories containing multiple `*.txt` files will no longer proceed automatically, even if one file was previously accepted by exact filename match.
- Update user-facing docs and configuration guidance to reflect the new transcript-source contract.
- Add behavior tests for `0/1/many .txt` discovery outcomes and ambiguity failures.

## Capabilities

### New Capabilities
- `transcript-source-discovery`: discover transcript input from `*.txt` files with strict confidence rules and consistent CLI behavior.

### Modified Capabilities
- None.

## Impact

- Affected code: `internal/summary/processor.go`, `internal/ui/filepicker.go`, `cmd/validate.go`, and related command messaging in `cmd/root.go`.
- Affected docs: `README.md`, `samples/settings.sample.yaml`, and CLI help text referencing fixed transcript filenames.
- Affected tests: summary behavior tests and command-level behavior tests for transcript discovery consistency.
