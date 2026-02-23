## Why

The main meeting summary includes detailed topic sections that are valuable for a shared drive but too verbose for Slack. Teams need a condensed, Slack-formatted version that surfaces the most actionable content (highlights, action items, risks) with a link back to the full summary. Generating both files in a single run eliminates manual copy-paste and ensures the extracted details are identical.

## What Changes

- Generate a second output file alongside the main summary with `-slack` appended to the filename (e.g., `2026-02-23-GISC-cadence-call-summary-slack.md`)
- The Slack mini summary contains the same extracted content (highlights, action items, risks) with zero deviation from the main summary — only the section order and included sections differ
- The Slack version omits detailed topic sections and reorders content for quick scanning:
  1. Title (bold-italic, same as main)
  2. Highlights
  3. Action Items
  4. Risks
  5. Meeting Recording
  6. Full Meeting Summary (placeholder link to main summary on shared drive)
- The Slack file is written only after the main summary passes validation — no Slack file if the main summary fails

## Capabilities

### New Capabilities
- `slack-mini-summary`: Generating a second, reordered Slack-friendly summary file from the validated main summary output, with configurable section order and an appended link to the full summary

### Modified Capabilities
- `safe-summary-generation-flow`: The generation flow now produces two output files instead of one; the Slack file is derived from the validated main summary and written in the same save step

## Impact

- **Code**: `internal/summary/processor.go` — new function to extract and reorder sections from the validated summary, new save path for the Slack file
- **Code**: `internal/app/service.go` — orchestration updated to include the Slack file in the run result
- **Output**: A second `.md` file appears in the meeting directory alongside the main summary
- **No breaking changes**: The main summary output, filename, and validation are unchanged
