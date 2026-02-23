## 1. Section Parser

- [x] 1.1 Implement `ParseSections(content string) map[string]string` in `internal/summary/slack_summary.go` that walks lines top-to-bottom and splits validated summary text into named sections (title, highlights, action-items, risks, meeting-recording, and topic sections) by detecting Slack markdown header patterns
- [x] 1.2 Add unit tests for `ParseSections` covering: summary with all sections present, summary without Risks, summary with multiple topic sections, edge case of empty input

## 2. Slack Summary Assembly

- [x] 2.1 Implement `BuildSlackSummary(sections map[string]string) string` in `internal/summary/slack_summary.go` that reassembles parsed sections in fixed order (Title, Highlights, Action Items, Risks, Meeting Recording) and appends the `*FULL MEETING SUMMARY*` section with placeholder link
- [x] 2.2 Add unit tests for `BuildSlackSummary` covering: all sections present, Risks omitted when absent, Full Meeting Summary always appended, topic sections excluded

## 3. Filename and Persistence

- [x] 3.1 Implement `GenerateSlackOutputFilename() (string, error)` in `internal/summary/slack_summary.go` that delegates to `GenerateOutputFilename()` and inserts `-slack` before `.md`
- [x] 3.2 Implement `SaveSlackSummary(content string) (string, error)` in `internal/summary/slack_summary.go` that writes the Slack mini summary to the meeting directory using the Slack filename
- [x] 3.3 Add unit tests for `GenerateSlackOutputFilename` covering: filename with date, filename without date

## 4. Orchestration

- [x] 4.1 Add `SlackOutputPath` and `SlackWarning` fields to `RunResult` in `internal/app/service.go`
- [x] 4.2 Update `Session.Run()` to call `ParseSections` → `BuildSlackSummary` → `SaveSlackSummary` after `SaveSummary` succeeds, treating Slack write failure as a non-fatal warning stored in `SlackWarning`
- [x] 4.3 Update CLI output in `cmd/root.go` to display the Slack file path on success or the warning on failure

## 5. Integration Testing

- [x] 5.1 Add an integration test in `internal/app/service_test.go` that verifies both main and Slack files are written to the meeting directory after a successful run
- [x] 5.2 Add an integration test that verifies the Slack file is not written when main summary validation fails
