## Context

The repository has a functioning CLI but mixes concerns across command handlers, domain orchestration, provider installation logic, and process-level side effects. CI and lint configuration also show drift from module/tooling baselines, and tests only cover a narrow subset of runtime behavior. A separate active change (`discover-transcript-from-txt`) defines transcript discovery behavior; modernization must preserve those requirements.

## Goals / Non-Goals

**Goals:**
- Separate CLI interaction concerns from summary workflow orchestration.
- Make runtime preflight provider-agnostic and driven by configured AI command.
- Remove process-wide side effects from summary generation and tighten output-write safety.
- Align CI/lint/local quality workflows with a single, deterministic Go/tooling baseline.
- Increase behavior-focused tests for runtime-preflight and generation-flow contracts.

**Non-Goals:**
- Replacing Cobra, Bubble Tea, or core UI libraries.
- Redesigning summary prompt content or transcript extraction semantics from the active transcript-discovery change.
- Introducing remote services, background daemons, or persistent databases.

## Decisions

1. Introduce a runtime service boundary
- Decision: Add an internal application service boundary that accepts an explicit run request (config, user identity, meeting directory, execution mode) and returns a structured result for CLI rendering.
- Rationale: This decouples business workflow from command/UI plumbing and improves testability.
- Alternatives considered:
- Keep logic in `cmd/root.go`: rejected due to complexity growth and low testability.
- Full plugin architecture now: rejected as excessive for current scope.

2. Provider-agnostic runtime preflight
- Decision: Runtime preflight will validate `config.ai.command` availability generically and return actionable errors based on configured command, while provider-specific installation remains in explicit install commands.
- Rationale: This matches configurable provider intent and avoids hard-coded Gemini assumptions in runtime path.
- Alternatives considered:
- Continue requiring Gemini in runtime regardless of config: rejected as contradictory behavior.
- Remove all install guidance: rejected because users still need clear next steps.

3. Eliminate process-wide working-directory mutation
- Decision: Stop using process-level `os.Chdir` in runtime flow; use explicit directory-aware execution (`exec.Cmd.Dir`) and explicit metadata/path injection.
- Rationale: Prevents global side effects and makes command execution safer and more composable.
- Alternatives considered:
- Keep `os.Chdir` with stricter defer handling: rejected because side effects remain global.

4. Validate before final write, retain diagnostics
- Decision: Validate generated summary content before writing final summary file; on validation failure, persist diagnostic raw output and return actionable error.
- Rationale: Prevents accidental empty/invalid final artifacts while preserving debugging evidence.
- Alternatives considered:
- Always write best-effort output: rejected due to risk of silent bad summaries.
- Fail without saving raw output: rejected due to poor diagnosability.

5. Single tooling baseline for CI and local workflow
- Decision: CI, lint, and local workflow commands will align with module Go version and deterministic quality gates; deprecated lint rules and placeholder build artifacts will be removed.
- Rationale: Reduces false confidence and avoids mismatch between local and CI behavior.
- Alternatives considered:
- Keep permissive/inconsistent tooling for convenience: rejected due to long-term maintenance cost.

6. Behavior-first testing strategy
- Decision: Expand behavior tests around runtime service, preflight behavior, and summary output safety using real filesystem fixtures and executable stubs where needed (no mocks).
- Rationale: Satisfies reliability goals and project testing rules while preserving realistic behavior validation.
- Alternatives considered:
- Continue mostly manual testing: rejected due to regression risk.

## Risks / Trade-offs

- [Service-boundary refactor introduces temporary churn] -> Mitigation: land refactor in small slices with behavior tests per slice.
- [Provider-agnostic preflight may reduce tailored guidance] -> Mitigation: provide generic guidance plus provider-specific hints when configured command is known.
- [Stricter output validation may increase user-visible failures initially] -> Mitigation: include clear failure reason and raw-output diagnostic location.
- [Tooling hardening can create short-term CI noise] -> Mitigation: update workflows and lint config in one coordinated change with clear migration notes.

## Migration Plan

1. Preserve active transcript-discovery semantics and avoid overlapping behavior regressions.
2. Introduce runtime service boundary and move orchestration out of command handlers.
3. Implement provider-agnostic preflight and update dependency guidance paths.
4. Remove process-wide cwd mutation and implement validate-before-write flow with diagnostics.
5. Align CI/lint/local tooling baselines and remove placeholder artifact behavior.
6. Add/expand behavior-focused tests for each migrated boundary.

## Open Questions

- Should the default `meetsum` command continue offering in-band dependency installation prompts, or should it only fail with guidance to explicit install commands?
- Where should invalid raw-output diagnostics be stored by default: meeting directory or log directory?

