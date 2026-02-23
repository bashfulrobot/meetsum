## Context

The current summary generation flow produces a single output file containing: Title, Topic Sections, Highlights, Action Items, and Meeting Recording. The AI generates all of this as one markdown document. After validation, it is written to disk.

Teams paste summaries into Slack but the topic sections are too verbose for that context. A condensed version is needed that keeps only the actionable sections (highlights, action items, risks) and adds a link to the full summary.

The validated main summary already contains all the content needed — the Slack mini summary is a pure text transformation, not a second AI call.

## Goals / Non-Goals

**Goals:**
- Generate a Slack mini summary file alongside the main summary in a single run
- Derive the Slack file entirely from the validated main summary (no second AI invocation)
- Preserve identical content in shared sections — zero drift between main and Slack versions
- Add a "Full Meeting Summary" section with a placeholder link for the shared drive URL

**Non-Goals:**
- Changing the main summary format or content
- Making the Slack section order configurable (hardcoded for now)
- Generating the Slack file independently of the main summary
- Populating the placeholder links with real URLs (manual step)

## Decisions

### 1. Parse sections from validated summary text

**Decision**: Implement a section parser that splits the validated summary into named sections by detecting header patterns (`*HIGHLIGHTS*`, `*ACTION ITEMS*`, `*RISKS*`, `*MEETING RECORDING*`, topic headers `_TOPIC_`).

**Rationale**: The summary is Slack-compatible markdown with predictable header patterns. Parsing the validated text guarantees identical content. Alternative — having the AI produce both formats — would double cost, introduce drift risk, and complicate validation.

**Approach**: Walk lines top-to-bottom, detect section boundaries by header patterns, accumulate lines into named sections. The title line (bold-italic `*_..._*`) is its own section.

### 2. Slack filename uses `-slack` suffix before extension

**Decision**: The Slack file uses the same base name with `-slack` inserted before `.md` (e.g., `2026-02-23-GISC-cadence-call-summary-slack.md`).

**Rationale**: Co-locates the two files in directory listings and makes the relationship obvious. The `-slack` suffix is appended by a `GenerateSlackOutputFilename()` method that delegates to `GenerateOutputFilename()` and inserts the suffix.

### 3. Slack file written inside the existing save step

**Decision**: `SaveSlackSummary()` is called immediately after `SaveSummary()` succeeds in `Session.Run()`. If the Slack file write fails, it is treated as a warning (non-fatal) — the main summary is already persisted.

**Rationale**: The main summary is the primary artifact. A failure in the secondary Slack file should not roll back or fail the entire run. The Slack path is returned in `RunResult` for CLI display.

**Alternative considered**: Making both writes atomic (both succeed or both fail). Rejected because the main summary is the source of truth; the Slack file can always be regenerated.

### 4. Fixed section order for Slack output

**Decision**: Hardcode the Slack section order as: Title, Highlights, Action Items, Risks, Meeting Recording, Full Meeting Summary. Sections not present in the source (e.g., no Risks) are simply omitted.

**Rationale**: This matches the user's required format. Making it configurable adds complexity with no current demand. If needed later, the ordered list can be extracted to config.

### 5. "Full Meeting Summary" section appended with placeholder link

**Decision**: Always append a `*FULL MEETING SUMMARY*` section with a placeholder link: `[Full Summary on Shared Drive](PLACEHOLDER_URL)`.

**Rationale**: Mirrors the existing `MEETING RECORDING` placeholder pattern. Users manually replace the URL after uploading to Google Drive.

## Risks / Trade-offs

- **Parsing fragility** → The section parser depends on consistent header formatting from the AI. Mitigation: the parser uses the same patterns the LLM instructions mandate (`*HEADER*` for bold sections, `_HEADER_` for topics, `*_TITLE_*` for the title). If the AI deviates, the main summary itself would look wrong — so the parser's assumptions align with existing quality expectations.
- **Missing sections** → If the AI doesn't generate a Risks section, the Slack output simply omits it. This is acceptable — the Slack file is a subset, not a superset.
- **Non-fatal Slack write failure** → If the Slack file can't be written (permissions, disk), the run still succeeds. The warning is surfaced to the user. Risk: user may not notice. Mitigation: CLI output explicitly shows both file paths or the warning.
