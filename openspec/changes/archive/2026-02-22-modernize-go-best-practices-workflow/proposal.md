## Why

The repository currently works but carries avoidable technical risk: runtime orchestration is tightly coupled in `cmd/root.go`, provider checks are partly hard-coded, CI/toolchain settings are inconsistent, and test coverage is minimal outside one package. Modernizing these foundations now will reduce regressions, improve maintainability, and make future feature work faster.

## What Changes

- Introduce clearer runtime boundaries so command handlers focus on CLI concerns while application services own summary workflow logic.
- Make runtime dependency checks provider-agnostic by using configured AI command behavior rather than hard-coded Gemini assumptions.
- Remove process-wide side effects from summary generation flow and tighten output validation before writing final summaries.
- Align CI, lint, and local quality workflows with a single Go/tooling baseline and remove stale/deprecated lint configuration.
- Expand behavior-focused tests across command/runtime boundaries and critical summary flow failure paths.
- Ensure modernization preserves behavior defined by `discover-transcript-from-txt` artifacts.

## Capabilities

### New Capabilities
- `provider-agnostic-runtime-preflight`: runtime preflight and user guidance are driven by configured AI command behavior, not provider-specific assumptions.
- `safe-summary-generation-flow`: summary generation avoids process-wide side effects, validates output before final write, and surfaces actionable failures.
- `workflow-quality-baseline`: CI, linting, and local quality commands follow a consistent Go/tooling contract and enforce meaningful build/test signals.

### Modified Capabilities
- None.

## Impact

- Affected code: `cmd/root.go`, `cmd/check.go`, `internal/summary/processor.go`, `internal/deps/gemini.go`, `config/config.go`, and supporting command/service wiring.
- Affected quality tooling: `.github/workflows/ci.yml`, `.golangci.yml`, and `justfile`.
- Affected tests: new behavior tests in `cmd` and `internal/summary` plus preflight/flow coverage.
- Change dependency: modernization implementation must preserve transcript-discovery behavior captured in `openspec/changes/discover-transcript-from-txt/`.

