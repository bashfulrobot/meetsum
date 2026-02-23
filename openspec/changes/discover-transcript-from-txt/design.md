## Context

`meetsum` currently assumes a fixed transcript filename in some code paths while other paths allow limited fallback behavior. This creates inconsistent outcomes between interactive directory selection, validation, and summary generation. The requested behavior is to treat any `*.txt` file in the meeting directory as a transcript candidate, while failing hard when selection confidence is low.

## Goals / Non-Goals

**Goals:**
- Establish one canonical transcript-discovery contract used by runtime generation, `meetsum validate`, and file-picker selection rules.
- Support transcript filenames that vary by customer/date folder structure by matching `*.txt` files.
- Enforce deterministic and safe confidence handling: exactly one candidate proceeds, zero candidates fail, and multiple candidates fail with listed candidates.
- Keep behavior portable across AI providers by keeping discovery independent from model-specific logic.

**Non-Goals:**
- Inferring transcript quality from content or file size.
- Recursive transcript discovery outside the selected meeting directory.
- Adding automatic tie-breaking heuristics (for example prioritizing `transcript.txt`) when multiple `*.txt` files exist.

## Decisions

1. Canonical candidate discovery rule
- Decision: Discover transcript candidates as files in the meeting directory whose extension is `.txt` (case-insensitive), then apply the `0/1/many` confidence policy.
- Rationale: This meets the requested flexibility while preventing silent mis-selection.
- Alternatives considered:
- Keep exact filename requirement: rejected because it fails common multi-customer naming patterns.
- Auto-pick the first/lexicographically smallest file: rejected because it can silently choose the wrong document.

2. Hard-fail on low confidence
- Decision: Fail for zero or multiple candidates, with actionable error messaging that includes discovered candidate names in deterministic order.
- Rationale: The user selected hard-fail behavior for low-confidence inference; this protects summary correctness.
- Alternatives considered:
- Warn-and-continue for interactive runs: rejected because it can still produce incorrect summaries.
- Prompt user to choose among candidates: rejected for now to keep deterministic CLI behavior and avoid adding a new interaction path.

3. Single shared discovery behavior across entry points
- Decision: Use one shared transcript-discovery behavior for processor validation, file picker acceptance rules, and `validate` command checks.
- Rationale: Prevents contract drift between commands and surfaces errors early and consistently.
- Alternatives considered:
- Keep command-specific rules: rejected because existing drift already causes confusion and contradictory outcomes.

4. Configuration compatibility treatment
- Decision: Keep existing transcript configuration fields for compatibility but treat them as non-authoritative for transcript discovery in this change; update docs accordingly.
- Rationale: Avoids unnecessary config-breaking churn while shifting behavior to the new canonical contract.
- Alternatives considered:
- Immediate config field removal: rejected because it introduces broader migration work beyond this behavioral change.

## Risks / Trade-offs

- [Directories with multiple text files now fail more often] -> Mitigation: clear ambiguity error listing candidates and documentation updates with expected folder hygiene.
- [Users may expect legacy `transcript.txt` preference] -> Mitigation: release notes and README explicitly describe the new confidence policy and breaking behavior.
- [No semantic quality check of chosen file] -> Mitigation: retain explicit hard-fail for ambiguous sets and consider future metadata-based transcript designation in a separate change.

## Migration Plan

1. Implement unified `*.txt` candidate discovery and confidence handling.
2. Update file picker and `validate` command to use the same contract.
3. Update README and sample settings/help text to remove fixed filename assumptions.
4. Add tests for `0/1/many` candidate outcomes and consistency across command paths.
5. Document breaking behavior in release notes.

## Open Questions

- Should a future change add an explicit transcript override flag or metadata field for teams that intentionally keep multiple `*.txt` files in one meeting folder?
