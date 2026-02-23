## Context

`meetsum` currently documents `ai.command` as a single command value, but runtime behavior already depends on parsing and executing that value. This makes provider flag usage unclear and fragile, especially when switching between AI CLIs with different argument contracts. We need an explicit, provider-agnostic way to configure optional CLI arguments while preserving existing no-flag and legacy command-string behavior.

## Goals / Non-Goals

**Goals:**
- Add first-class configuration for optional AI CLI arguments that works across providers.
- Preserve current behavior when no arguments are configured.
- Preserve backward compatibility for existing configs that embed flags in `ai.command`.
- Keep runtime preflight, dependency checks, and execution aligned to one resolved invocation contract.
- Add behavior-focused tests using real executable scripts (no mocks).

**Non-Goals:**
- Introducing provider-specific configuration schemas (for example separate `openai` or `anthropic` config blocks).
- Supporting shell expression parsing, interpolation, or `bash -c` style execution.
- Implementing automatic provider-specific flag validation.
- Redesigning install/docs subcommands beyond what is needed for accurate invocation guidance.

## Decisions

1. Add explicit `ai.args` list configuration while retaining `ai.command`
- Decision: Extend config with `ai.args` (`[]string`) and keep `ai.command` as the primary executable field.
- Rationale: A string list is portable and avoids shell-quoting ambiguity. It supports all providers because flags remain opaque pass-through tokens.
- Alternatives considered:
- Continue with a single free-form `ai.command` string only: rejected because quoting and portability are brittle and undocumented for multi-argument usage.
- Add provider-specific model/flag keys: rejected because it creates coupling to specific CLIs and does not scale across tools.

2. Resolve one canonical invocation tuple before preflight and execution
- Decision: Introduce a pure resolver in `internal/ai` that returns `(executable, args)` from config with compatibility rules:
- If `ai.args` is non-empty, use `ai.command` as executable and append `ai.args` as-is.
- If `ai.args` is empty, support legacy inline flags by parsing `ai.command` into executable + args with current whitespace behavior.
- If `ai.args` is non-empty and `ai.command` appears to contain inline args, fail validation with actionable guidance to avoid ambiguity.
- Rationale: One resolver eliminates drift between preflight checks and runtime execution while keeping old configs functional.
- Alternatives considered:
- Hard break and require immediate migration to `ai.args`: rejected to avoid unnecessary config breakage.
- Execute raw command via shell: rejected for safety, portability, and quoting unpredictability.

3. Keep preflight and dependency checks provider-agnostic, keyed on resolved executable
- Decision: `Preflight`, `meetsum check`, and runtime guidance must use the resolved executable for PATH checks; provider-specific Gemini checks continue only when the resolved executable target is Gemini.
- Rationale: This preserves provider-agnostic behavior while keeping existing optional Gemini guidance when relevant.
- Alternatives considered:
- Remove Gemini-specific guidance entirely: rejected as out of scope for this change.
- Keep checks keyed on raw config text: rejected because it breaks when args/path variants are configured.

4. Treat argument display as an observability concern with safe defaults
- Decision: User-facing status output should show command and configured-args presence without requiring full argument echo in all paths.
- Rationale: Some providers accept sensitive flags; minimizing raw arg echo reduces accidental leakage while still exposing enough diagnostic context.
- Alternatives considered:
- Always print full resolved args: rejected due to potential credential leakage.
- Never show argument context: rejected because it weakens troubleshooting.

5. Validate behavior via end-to-end invocation tests with real scripts
- Decision: Add tests that execute temporary fake CLIs and assert observed argv behavior for:
- No args configured.
- `ai.args` configured.
- Legacy inline args in `ai.command`.
- Ambiguous mixed config rejected.
- Rationale: Behavior tests directly verify external process invocation contracts and align with repository test policy.
- Alternatives considered:
- Parser-only unit tests: rejected because they do not prove execution behavior.

## Risks / Trade-offs

- [Legacy inline parsing still uses whitespace tokenization] -> Mitigation: document that complex quoting should use `ai.args`; preserve legacy behavior only for compatibility.
- [Mixed config ambiguity (`ai.command` with inline args plus `ai.args`)] -> Mitigation: fail fast with actionable config error.
- [Reduced argument echo may hide debugging detail] -> Mitigation: provide explicit trace-level output guidance and safe summaries (for example executable + arg count).

## Migration Plan

1. Extend config schema and defaults with `ai.args`.
2. Implement canonical invocation resolver and validation rules in `internal/ai`.
3. Update preflight, runtime execution, and dependency-check command paths to use the canonical resolver.
4. Update README and sample settings with provider-agnostic examples using `ai.command` + `ai.args`.
5. Add behavior tests for invocation outcomes and compatibility rules.

## Open Questions

- Should environment variable overrides support list-style argument values directly (for example CSV or repeated env keys), or remain YAML-first for `ai.args`?
- Do we want a standardized redaction policy for argument keys that are commonly sensitive (for example flags containing `key`, `token`, or `secret`)?

