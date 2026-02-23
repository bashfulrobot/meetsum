## Why

`meetsum` currently exposes only a single configured AI command, which limits users who need provider-specific flags such as model selection. We need a provider-agnostic way to configure optional CLI arguments so teams can switch tools without code changes or fragile command-string workarounds.

## What Changes

- Add configuration support for optional AI CLI arguments in `settings.yaml` alongside `ai.command`.
- Define runtime behavior so meetsum invokes the configured executable with configured arguments and continues to pass prompt content through stdin.
- Preserve current behavior when no arguments are configured (empty/default args).
- Keep dependency preflight and `meetsum check` provider-agnostic by validating configured command availability while surfacing configured invocation details clearly.
- Document configuration examples for different AI CLIs and argument styles.
- Add behavior-focused tests covering no-arg invocation, configured-arg invocation, and cross-provider command switching.

## Capabilities

### New Capabilities
- `ai-cli-invocation-config`: configure command and optional arguments for AI CLI execution in a provider-agnostic, backward-compatible way.

### Modified Capabilities
- None.

## Impact

- Affected code: `config/config.go`, `internal/ai/command.go`, `internal/summary/processor.go`, `internal/app/service.go`, and command surfaces in `cmd/check.go`/`cmd/root.go`.
- Affected docs/config samples: `README.md` and `samples/settings.sample.yaml`.
- Affected tests: config loading tests, AI command parsing/invocation tests, and runtime behavior tests for provider switching with and without args.

