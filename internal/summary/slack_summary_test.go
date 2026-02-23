package summary

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testSummaryAllSections = `*_2026-02-23 ACME CADENCE CALL SUMMARY_*

_PRODUCT ROADMAP_

We discussed the upcoming product roadmap for Q2.

_PARTNERSHIP UPDATES_

The partnership with XYZ Corp is progressing well.

*HIGHLIGHTS*

- Acme is moving forward with the enterprise plan.
- The new integration will enable automated data sync.

*ACTION ITEMS*

- John Doe: Send the proposal document by Friday.
- Jane Smith: Schedule a follow-up call for next week.

*RISKS*

- Timeline Risk: The Q2 launch may be delayed.

*MEETING RECORDING*

- [Meeting Recording](PLACEHOLDER_URL)`

const testSummarySlashTopic = `*_2026-02-24 GISC CADENCE CALL SUMMARY_*

_CURRENT STATE/USE CASES_

The GISC team has made good progress in their proof-of-concept.

*HIGHLIGHTS*

- Validated Kong proxy capabilities with existing MCPs on AWS.

*ACTION ITEMS*

- Joe Cesario: Send an invite for the Professional Services scoping call.

*RISKS*

- Timeline Risk: The target of launching by end of April is aggressive.

*MEETING RECORDING*

- [Clari Recording](PLACEHOLDER_URL)`

const testSummaryNoRisks = `*_2026-02-23 ACME CADENCE CALL SUMMARY_*

_PRODUCT ROADMAP_

We discussed the upcoming product roadmap.

*HIGHLIGHTS*

- Key insight from the meeting.

*ACTION ITEMS*

- John Doe: Complete the analysis.

*MEETING RECORDING*

- [Meeting Recording](PLACEHOLDER_URL)`

func TestParseSections(t *testing.T) {
	t.Run("all sections present", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)

		expectedKeys := []string{
			"title", "highlights", "action-items", "risks",
			"meeting-recording", "topic:PRODUCT ROADMAP", "topic:PARTNERSHIP UPDATES",
		}
		for _, key := range expectedKeys {
			if _, ok := sections[key]; !ok {
				t.Errorf("expected section %q not found", key)
			}
		}

		if !strings.HasPrefix(sections["title"], "*_2026-02-23 ACME CADENCE CALL SUMMARY_*") {
			t.Errorf("title content mismatch: %s", sections["title"])
		}
		if !strings.Contains(sections["highlights"], "enterprise plan") {
			t.Errorf("highlights content mismatch: %s", sections["highlights"])
		}
		if !strings.Contains(sections["action-items"], "John Doe") {
			t.Errorf("action-items content mismatch: %s", sections["action-items"])
		}
		if !strings.Contains(sections["risks"], "Timeline Risk") {
			t.Errorf("risks content mismatch: %s", sections["risks"])
		}
		if !strings.Contains(sections["meeting-recording"], "PLACEHOLDER_URL") {
			t.Errorf("meeting-recording content mismatch: %s", sections["meeting-recording"])
		}
	})

	t.Run("summary without risks", func(t *testing.T) {
		sections := ParseSections(testSummaryNoRisks)

		if _, ok := sections["risks"]; ok {
			t.Error("expected no risks section")
		}
		if _, ok := sections["title"]; !ok {
			t.Error("expected title section")
		}
		if _, ok := sections["highlights"]; !ok {
			t.Error("expected highlights section")
		}
		if _, ok := sections["action-items"]; !ok {
			t.Error("expected action-items section")
		}
	})

	t.Run("multiple topic sections", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)

		if _, ok := sections["topic:PRODUCT ROADMAP"]; !ok {
			t.Error("expected PRODUCT ROADMAP topic")
		}
		if _, ok := sections["topic:PARTNERSHIP UPDATES"]; !ok {
			t.Error("expected PARTNERSHIP UPDATES topic")
		}
	})

	t.Run("topic with slash in name", func(t *testing.T) {
		sections := ParseSections(testSummarySlashTopic)

		if _, ok := sections["topic:CURRENT STATE/USE CASES"]; !ok {
			t.Error("expected topic:CURRENT STATE/USE CASES section")
		}

		// Title should only contain the title line, not the topic content
		if strings.Contains(sections["title"], "proof-of-concept") {
			t.Error("topic content should not be lumped into title section")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		sections := ParseSections("")
		if len(sections) != 0 {
			t.Errorf("expected empty map for empty input, got %d sections", len(sections))
		}
	})
}

func TestBuildSlackSummary(t *testing.T) {
	t.Run("all sections present in correct order", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)
		result := BuildSlackSummary(sections)

		titleIdx := strings.Index(result, "*_2026-02-23 ACME CADENCE CALL SUMMARY_*")
		highlightsIdx := strings.Index(result, "*HIGHLIGHTS*")
		actionIdx := strings.Index(result, "*ACTION ITEMS*")
		risksIdx := strings.Index(result, "*RISKS*")
		recordingIdx := strings.Index(result, "*MEETING RECORDING*")
		fullSummaryIdx := strings.Index(result, "*FULL MEETING SUMMARY*")

		if titleIdx < 0 || highlightsIdx < 0 || actionIdx < 0 || risksIdx < 0 || recordingIdx < 0 || fullSummaryIdx < 0 {
			t.Fatalf("missing sections in result:\n%s", result)
		}

		if !(titleIdx < highlightsIdx && highlightsIdx < actionIdx && actionIdx < risksIdx && risksIdx < recordingIdx && recordingIdx < fullSummaryIdx) {
			t.Error("sections out of order")
		}
	})

	t.Run("risks omitted when absent", func(t *testing.T) {
		sections := ParseSections(testSummaryNoRisks)
		result := BuildSlackSummary(sections)

		if strings.Contains(result, "*RISKS*") {
			t.Error("expected no RISKS section in output")
		}
		if !strings.Contains(result, "*FULL MEETING SUMMARY*") {
			t.Error("Full Meeting Summary should always be present")
		}
	})

	t.Run("full meeting summary always appended", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)
		result := BuildSlackSummary(sections)

		if !strings.Contains(result, "*FULL MEETING SUMMARY*") {
			t.Error("missing Full Meeting Summary section")
		}
		if !strings.Contains(result, "[Full Summary on Shared Drive](PLACEHOLDER_URL)") {
			t.Error("missing placeholder link")
		}
	})

	t.Run("topic sections excluded", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)
		result := BuildSlackSummary(sections)

		if strings.Contains(result, "_PRODUCT ROADMAP_") {
			t.Error("topic section should be excluded from Slack output")
		}
		if strings.Contains(result, "_PARTNERSHIP UPDATES_") {
			t.Error("topic section should be excluded from Slack output")
		}
		if strings.Contains(result, "We discussed the upcoming product roadmap") {
			t.Error("topic content should be excluded from Slack output")
		}
	})

	t.Run("slash topic excluded from slack output", func(t *testing.T) {
		sections := ParseSections(testSummarySlashTopic)
		result := BuildSlackSummary(sections)

		if strings.Contains(result, "_CURRENT STATE/USE CASES_") {
			t.Error("topic with slash should be excluded from Slack output")
		}
		if strings.Contains(result, "proof-of-concept") {
			t.Error("topic content should be excluded from Slack output")
		}
	})

	t.Run("content identity with source sections", func(t *testing.T) {
		sections := ParseSections(testSummaryAllSections)
		result := BuildSlackSummary(sections)

		slackSections := ParseSections(result)

		for _, key := range []string{"title", "highlights", "action-items", "risks", "meeting-recording"} {
			original, hasOriginal := sections[key]
			slack, hasSlack := slackSections[key]

			if hasOriginal != hasSlack {
				t.Errorf("section %q presence mismatch: original=%v slack=%v", key, hasOriginal, hasSlack)
				continue
			}
			if hasOriginal && original != slack {
				t.Errorf("section %q content drift:\noriginal: %q\nslack:    %q", key, original, slack)
			}
		}
	})
}

func TestGenerateSlackOutputFilename(t *testing.T) {
	t.Run("filename with date", func(t *testing.T) {
		baseDir := t.TempDir()
		meetingDir := filepath.Join(baseDir, "Customers", "GISC", "2026-02-23")
		if err := os.MkdirAll(meetingDir, 0755); err != nil {
			t.Fatalf("failed to create meeting dir: %v", err)
		}

		processor := newTestProcessor(t, meetingDir)
		filename, err := processor.GenerateSlackOutputFilename()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if filename != "2026-02-23-GISC-cadence-call-summary-slack.md" {
			t.Fatalf("expected 2026-02-23-GISC-cadence-call-summary-slack.md, got %s", filename)
		}
	})

	t.Run("filename without date", func(t *testing.T) {
		baseDir := t.TempDir()
		meetingDir := filepath.Join(baseDir, "Customers", "GISC", "no-date-here")
		if err := os.MkdirAll(meetingDir, 0755); err != nil {
			t.Fatalf("failed to create meeting dir: %v", err)
		}

		processor := newTestProcessor(t, meetingDir)
		filename, err := processor.GenerateSlackOutputFilename()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if filename != "GISC-cadence-call-summary-slack.md" {
			t.Fatalf("expected GISC-cadence-call-summary-slack.md, got %s", filename)
		}
	})
}
