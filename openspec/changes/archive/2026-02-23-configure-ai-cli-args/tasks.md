## 1. Configuration Model

- [x] 1.1 Add `ai.args` (`[]string`) to configuration schema and defaults in `config/config.go`.
- [x] 1.2 Update configuration display output in `cmd/config.go` to show AI argument configuration safely (for example arg count or non-sensitive summary).
- [x] 1.3 Update `samples/settings.sample.yaml` with provider-agnostic examples showing `ai.command` with optional `ai.args`.

## 2. Invocation Resolution and Validation

- [x] 2.1 Add a canonical resolver in `internal/ai` that returns resolved executable and args from `ai.command` + `ai.args`.
- [x] 2.2 Preserve legacy inline parsing when `ai.args` is empty, and add explicit validation failure for mixed inline args plus non-empty `ai.args`.
- [x] 2.3 Update command availability checks to validate the resolved executable from the canonical resolver.

## 3. Runtime and Check Command Integration

- [x] 3.1 Update runtime preflight/service flow to use the canonical invocation resolver and propagate actionable errors.
- [x] 3.2 Update summary processor execution to use resolved executable and arguments while preserving stdin prompt flow and stdout/stderr capture.
- [x] 3.3 Update `meetsum check` and preflight guidance output to remain provider-agnostic and use resolved executable behavior for dependency messaging.

## 4. Documentation and Behavior Tests

- [x] 4.1 Update `README.md` configuration docs to describe optional `ai.args`, no-arg behavior, and provider/tool switching guidance.
- [x] 4.2 Add behavior-focused tests for no-arg invocation and configured-arg invocation using real temporary executables.
- [x] 4.3 Add behavior-focused tests for legacy inline `ai.command` compatibility and mixed-source configuration rejection.
