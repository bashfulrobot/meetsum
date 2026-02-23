## 1. Runtime Boundary Refactor

- [x] 1.1 Introduce an internal runtime service boundary that moves summary orchestration logic out of `cmd/root.go`.
- [x] 1.2 Refactor command handlers to keep CLI interaction/rendering concerns separate from runtime workflow execution.
- [x] 1.3 Add behavior tests validating service-level run outcomes for success and failure paths.

## 2. Provider-Agnostic Preflight and Safe Generation Flow

- [x] 2.1 Replace hard-coded runtime provider checks with configured-command preflight behavior.
- [x] 2.2 Remove process-wide `os.Chdir` usage from summary generation flow and use explicit execution directory context.
- [x] 2.3 Implement validate-before-write summary output flow with preserved diagnostic raw output on validation failure.
- [x] 2.4 Add behavior tests for configured-command preflight outcomes and summary validation-before-write failures.
- [x] 2.5 Verify modernization changes preserve transcript discovery behavior defined in `discover-transcript-from-txt`.

## 3. CI, Lint, and Local Quality Baseline

- [x] 3.1 Align CI Go toolchain version with module baseline in `go.mod`.
- [x] 3.2 Update build workflow to publish only real built artifacts for declared targets.
- [x] 3.3 Pin workflow action references to stable version tags and remove mutable branch references.
- [x] 3.4 Remove deprecated lint entries and stabilize lint configuration behavior.
- [x] 3.5 Provide a deterministic local quality command path (format, vet, tests) and document expected usage.
