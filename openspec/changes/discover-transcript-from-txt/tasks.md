## 1. Transcript Discovery Contract

- [x] 1.1 Add shared transcript-candidate discovery that finds meeting-directory files with `.txt` extension (case-insensitive).
- [x] 1.2 Implement deterministic `0/1/many` confidence handling (select one, fail on none, fail on many with candidate list).
- [x] 1.3 Update summary preprocessing flow to use the shared discovery result as the transcript source.

## 2. Command Path Consistency

- [x] 2.1 Update file picker directory acceptance to enforce the same transcript discovery contract.
- [x] 2.2 Update `meetsum validate` meeting-directory checks to enforce the same transcript discovery contract.
- [x] 2.3 Ensure transcript-related error messaging is actionable and consistent across runtime, picker, and validate paths.

## 3. Tests and Documentation

- [x] 3.1 Add behavior tests for transcript discovery with `0/1/many .txt` files, including deterministic ambiguity reporting.
- [x] 3.2 Add command-level behavior tests (or equivalent coverage) proving consistent pass/fail outcomes across generation and validation paths.
- [x] 3.3 Update README and sample configuration/help text to document `.txt` transcript discovery and ambiguity hard-fail behavior.
