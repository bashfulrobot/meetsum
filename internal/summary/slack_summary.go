package summary

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitfield/script"
)

var (
	boldItalicHeaderRe = regexp.MustCompile(`^\*_.*_\*$`)
	boldHeaderRe       = regexp.MustCompile(`^\*([A-Z][A-Z ]*)\*$`)
	italicHeaderRe     = regexp.MustCompile(`^_([A-Z][A-Z0-9 /&-]*?)_$`)
)

// ParseSections splits validated summary text into named sections by detecting
// Slack markdown header patterns. Returns a map keyed by normalized section name.
// Bold-italic headers (*_..._*) map to "title".
// Bold headers (*SECTION*) map to lowercase hyphenated keys (e.g., "action-items").
// Italic-only headers (_TOPIC_) map to "topic:NAME" keys.
func ParseSections(content string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentKey string
	var currentLines []string

	flush := func() {
		if currentKey != "" {
			text := strings.Join(currentLines, "\n")
			sections[currentKey] = strings.TrimRight(text, " \t\n")
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if boldItalicHeaderRe.MatchString(trimmed) {
			flush()
			currentKey = "title"
			currentLines = []string{line}
			continue
		}

		if boldHeaderRe.MatchString(trimmed) {
			flush()
			name := boldHeaderRe.FindStringSubmatch(trimmed)[1]
			currentKey = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			currentLines = []string{line}
			continue
		}

		if italicHeaderRe.MatchString(trimmed) {
			flush()
			name := italicHeaderRe.FindStringSubmatch(trimmed)[1]
			currentKey = "topic:" + name
			currentLines = []string{line}
			continue
		}

		if currentKey != "" {
			currentLines = append(currentLines, line)
		}
	}

	flush()
	return sections
}

// slackSectionOrder defines the fixed section key order for Slack mini summaries.
var slackSectionOrder = []string{
	"title",
	"highlights",
	"action-items",
	"risks",
	"meeting-recording",
}

// BuildSlackSummary reassembles parsed sections in the fixed Slack order and
// appends a Full Meeting Summary section with a placeholder link.
func BuildSlackSummary(sections map[string]string) string {
	var parts []string

	for _, key := range slackSectionOrder {
		if content, ok := sections[key]; ok {
			parts = append(parts, content)
		}
	}

	parts = append(parts, "*FULL MEETING SUMMARY*\n\n* [Full Summary on Shared Drive](PLACEHOLDER_URL)")

	return strings.Join(parts, "\n\n")
}

// GenerateSlackOutputFilename creates the Slack output filename by inserting
// -slack before the .md extension of the main summary filename.
func (p *Processor) GenerateSlackOutputFilename() (string, error) {
	mainFilename, err := p.GenerateOutputFilename()
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(mainFilename)
	base := strings.TrimSuffix(mainFilename, ext)
	return base + "-slack" + ext, nil
}

// SaveSlackSummary writes the Slack mini summary to the meeting directory.
func (p *Processor) SaveSlackSummary(content string) (string, error) {
	filename, err := p.GenerateSlackOutputFilename()
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(p.meetingDir, filename)

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	n, err := script.Echo(content).WriteFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to save slack summary: %w", err)
	}
	_ = n

	return outputPath, nil
}
