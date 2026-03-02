package summary

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bashfulrobot/meetsum/internal/ai"
	"github.com/bitfield/script"
	"github.com/charmbracelet/log"
)

type Processor struct {
	config         *config.Config
	logger         *log.Logger
	userName       string
	meetingDir     string
	transcriptPath string
}

// GeneratedSummaryOutput captures both cleaned and raw AI output.
type GeneratedSummaryOutput struct {
	Cleaned string
	Raw     string
}

// NewProcessor creates a new summary processor
func NewProcessor(cfg *config.Config, logger *log.Logger) *Processor {
	return &Processor{
		config: cfg,
		logger: logger,
	}
}

// SetUserName sets the user name for first-person perspective
func (p *Processor) SetUserName(name string) {
	p.userName = name
}

// SetMeetingDir sets the meeting directory
func (p *Processor) SetMeetingDir(dir string) {
	p.meetingDir = dir
}

// FindTranscriptFile resolves transcript source with the 0/1/many .txt contract.
func (p *Processor) FindTranscriptFile() (string, error) {
	return FindSingleTranscriptCandidate(p.meetingDir)
}

// ValidateRequiredFiles checks if all required files exist
func (p *Processor) ValidateRequiredFiles() error {
	// Find and validate transcript file
	transcriptPath, err := p.FindTranscriptFile()
	if err != nil {
		return err
	}
	p.transcriptPath = transcriptPath

	// Check instructions file
	instructionsPath := p.config.GetInstructionsPath()
	if _, err := os.Stat(instructionsPath); os.IsNotExist(err) {
		return fmt.Errorf("instructions file not found at %s", instructionsPath)
	}

	return nil
}

// GetOptionalFiles returns list of optional files that exist
func (p *Processor) GetOptionalFiles() []string {
	var files []string
	povPath := p.config.GetPovInputPath(p.meetingDir)
	if _, err := os.Stat(povPath); err == nil {
		files = append(files, "📝 pov-input.md")
	}
	return files
}

// LoadInstructions reads the LLM instructions file
func (p *Processor) LoadInstructions() (string, error) {
	instructionsPath := p.config.GetInstructionsPath()
	content, err := script.File(instructionsPath).String()
	if err != nil {
		return "", fmt.Errorf("failed to load instructions: %w", err)
	}
	return content, nil
}

// LoadTranscript reads the transcript file
func (p *Processor) LoadTranscript() (string, error) {
	if p.transcriptPath == "" {
		return "", fmt.Errorf("transcript path not set; call ValidateRequiredFiles first")
	}
	content, err := script.File(p.transcriptPath).String()
	if err != nil {
		return "", fmt.Errorf("failed to load transcript: %w", err)
	}
	return content, nil
}

// TranscriptPath returns the transcript path selected during validation.
func (p *Processor) TranscriptPath() string {
	return p.transcriptPath
}

// LoadContext reads the optional POV input file if it exists
func (p *Processor) LoadContext() (string, error) {
	povPath := p.config.GetPovInputPath(p.meetingDir)
	if _, err := os.Stat(povPath); os.IsNotExist(err) {
		return "", nil // No context file, return empty string
	}

	content, err := script.File(povPath).String()
	if err != nil {
		return "", fmt.Errorf("failed to load context: %w", err)
	}

	return fmt.Sprintf("CONTEXT GUIDE:\n%s", content), nil
}

// ExtractCustomerName extracts customer name from the meeting directory path.
// Returns (proper-case name, UPPER-CASE name).
func (p *Processor) ExtractCustomerName() (string, string) {
	// Extract customer name from path like /home/dustin/Documents/Kong/Customers/CustomerName/date
	customerNameRaw := ""
	if strings.Contains(p.meetingDir, "/Customers/") {
		parts := strings.Split(p.meetingDir, "/Customers/")
		if len(parts) > 1 {
			customerParts := strings.Split(parts[1], "/")
			if len(customerParts) > 0 {
				customerNameRaw = customerParts[0]
			}
		}
	}

	// Fallback to parent directory name
	if customerNameRaw == "" {
		parentDir := filepath.Dir(p.meetingDir)
		customerNameRaw = filepath.Base(parentDir)
	}

	return customerNameRaw, strings.ToUpper(customerNameRaw)
}

// ExtractDateFromPath extracts date from folder path like "2025-09-24"
func (p *Processor) ExtractDateFromPath() string {
	// Look for date pattern YYYY-MM-DD in the path
	parts := strings.Split(p.meetingDir, "/")
	for _, part := range parts {
		// Check if part matches YYYY-MM-DD format
		if len(part) == 10 && part[4] == '-' && part[7] == '-' {
			// Validate it's actually a date
			if _, err := time.Parse("2006-01-02", part); err == nil {
				return part
			}
		}
	}
	return "" // No date found in path
}

// GenerateOutputFilename creates the output filename
func (p *Processor) GenerateOutputFilename() (string, error) {
	name, _ := p.ExtractCustomerName()

	date := p.ExtractDateFromPath()
	if date == "" {
		return fmt.Sprintf("%s-cadence-call-summary.md", name), nil
	}
	return fmt.Sprintf("%s-%s-cadence-call-summary.md", date, name), nil
}

// GenerateSummary processes the meeting and generates a cleaned summary.
func (p *Processor) GenerateSummary() (string, error) {
	output, err := p.GenerateSummaryOutput()
	if err != nil {
		return "", err
	}
	return output.Cleaned, nil
}

// GenerateSummaryOutput processes the meeting and returns cleaned + raw output.
func (p *Processor) GenerateSummaryOutput() (GeneratedSummaryOutput, error) {
	// Load all required content
	instructions, err := p.LoadInstructions()
	if err != nil {
		return GeneratedSummaryOutput{}, err
	}

	transcript, err := p.LoadTranscript()
	if err != nil {
		return GeneratedSummaryOutput{}, err
	}

	context, err := p.LoadContext()
	if err != nil {
		return GeneratedSummaryOutput{}, err
	}

	// Extract date and customer info for the prompt
	customerNameProper, customerNameUpper := p.ExtractCustomerName()

	date := p.ExtractDateFromPath()
	titleDate := "UNDATED"
	if date != "" {
		titleDate = date
	}

	transcriptFile := filepath.Base(p.transcriptPath)

	// Prepare the prompt
	prompt := fmt.Sprintf(`%s

Process the transcript in %s and generate a structured meeting summary following the provided instructions. Write the summary from %s's first-person perspective.

The meeting date should be: %s
The customer name should be: %s (uppercase: %s)

IMPORTANT OUTPUT INSTRUCTIONS:
- Output ONLY the summary content directly with its Slack-compatible markdown formatting intact.
- Do NOT wrap the output in triple-backtick code fences.
- Do NOT attempt to save, write, or create any files.
- Do NOT include any preamble, postamble, or conversational text such as "Here is the summary" or "Is there anything else".
- The ENTIRE output must be the summary itself and nothing else.

TRANSCRIPT:
%s

%s`, instructions, transcriptFile, p.userName, titleDate, customerNameProper, customerNameUpper, transcript, context)

	// Execute AI command with separate stdout/stderr capture
	result, stderr, err := p.executeAICommand(prompt)
	if err != nil {
		p.logCommandError(stderr, err)
		return GeneratedSummaryOutput{}, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Clean the AI output to extract only the markdown content
	cleanedResult := p.cleanAIOutput(result)
	return GeneratedSummaryOutput{
		Cleaned: cleanedResult,
		Raw:     result,
	}, nil
}

// executeAICommand runs the AI command and captures stdout/stderr separately
func (p *Processor) executeAICommand(prompt string) (stdout string, stderr string, err error) {
	command, args, err := ai.ResolveConfiguredInvocation(p.config.AI.Command, p.config.AI.Args)
	if err != nil {
		return "", "", err
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = p.meetingDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

// logCommandError logs the command error with full context
func (p *Processor) logCommandError(stderr string, err error) {
	if p.logger == nil {
		return
	}

	p.logger.Error("AI command failed",
		"command", p.config.AI.Command,
		"error", err.Error(),
		"stderr", strings.TrimSpace(stderr),
		"meeting_dir", p.meetingDir,
	)
}

// cleanAIOutput removes error messages and extracts only the markdown content
func (p *Processor) cleanAIOutput(output string) string {
	lines := strings.Split(output, "\n")
	var markdownLines []string
	inMarkdownBlock := false
	foundMarkdownStart := false

	for _, line := range lines {
		// Skip error messages and AI chatter
		if isAIChatterLine(line) {
			continue
		}

		// Look for markdown code blocks
		if strings.HasPrefix(line, "```markdown") {
			inMarkdownBlock = true
			foundMarkdownStart = true
			continue
		}

		if strings.HasPrefix(line, "```") && inMarkdownBlock {
			inMarkdownBlock = false
			break // End of markdown block
		}

		// If we're in a markdown block, collect the content
		if inMarkdownBlock {
			markdownLines = append(markdownLines, line)
		}

		// If no markdown block found, but we see content that looks like a summary title
		if !inMarkdownBlock && !foundMarkdownStart && (strings.Contains(line, "_SUMMARY_") || strings.HasPrefix(line, "*_")) {
			// This looks like the start of the actual summary, collect everything from here
			markdownLines = append(markdownLines, line)
			foundMarkdownStart = true
		} else if foundMarkdownStart && !inMarkdownBlock {
			// We're collecting content after finding the summary start
			markdownLines = append(markdownLines, line)
		}
	}

	// If we found markdown content, use it
	if len(markdownLines) > 0 {
		return strings.TrimSpace(strings.Join(markdownLines, "\n"))
	}

	// Fallback: strip all chatter lines from the raw output
	var fallbackLines []string
	for _, line := range lines {
		if isAIChatterLine(line) {
			continue
		}
		fallbackLines = append(fallbackLines, line)
	}

	cleanedOutput := strings.Join(fallbackLines, "\n")

	// Remove any remaining markdown code block markers
	cleanedOutput = strings.ReplaceAll(cleanedOutput, "```markdown", "")
	cleanedOutput = strings.ReplaceAll(cleanedOutput, "```", "")

	return strings.TrimSpace(cleanedOutput)
}

// aiChatterPatterns are substrings that indicate a line is AI conversational
// noise rather than actual summary content. Patterns must be specific enough
// to avoid false positives against legitimate meeting summary text (e.g.
// action items like "Let me know if you need the report").
var aiChatterPatterns = []string{
	// Gemini CLI stderr/tool noise
	"Loaded cached credentials",
	"Error executing tool",
	"Tool \"write_file\" not found",
	// Conversational preamble/postamble from the AI
	"I was unable to create the file",
	"Here is the content for the meeting summary",
	"You can save it as",
	"I have generated the meeting summary",
	"I've generated the meeting summary",
	"Is there anything else I can help you with",
	"Is there anything else you would like",
	"Is there anything else you'd like",
	"I have saved the summary",
	"I've saved the summary",
	"I have created the file",
	"I've created the file",
	"I have written the summary",
	"I've written the summary",
	"Here's the meeting summary",
	"Here is the meeting summary",
	"Below is the meeting summary",
	"The summary has been saved",
	"The summary has been written",
	"and saved it to",
}

// isAIChatterLine returns true if a line matches known AI conversational noise.
func isAIChatterLine(line string) bool {
	for _, pattern := range aiChatterPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	return false
}

// ValidateSummaryContent ensures generated summary output is safe to persist.
func (p *Processor) ValidateSummaryContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("generated summary output is empty after cleaning")
	}

	// Reject output that is purely AI conversational chatter (no actual summary).
	lines := strings.Split(strings.TrimSpace(content), "\n")
	nonEmpty := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmpty++
		}
	}
	if nonEmpty <= 3 {
		// Very short output is suspicious; check for chatter patterns.
		for _, pattern := range aiChatterPatterns {
			if strings.Contains(content, pattern) {
				return fmt.Errorf("generated summary appears to contain AI conversational text rather than actual summary content")
			}
		}
	}

	return nil
}

// SaveRawOutputDiagnostics stores raw AI output when summary validation fails.
func (p *Processor) SaveRawOutputDiagnostics(rawOutput string) (string, error) {
	diagnosticPath := filepath.Join(p.meetingDir, "summary-raw-output.txt")
	content := strings.TrimSpace(rawOutput) + "\n"
	if err := os.WriteFile(diagnosticPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to save raw output diagnostics: %w", err)
	}

	return diagnosticPath, nil
}

// writeContentFile writes content to path, ensuring a trailing newline.
func writeContentFile(content, path string) error {
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	_, err := script.Echo(content).WriteFile(path)
	return err
}

// SaveSummary saves the generated summary to a file
func (p *Processor) SaveSummary(content string) (string, error) {
	filename, err := p.GenerateOutputFilename()
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(p.meetingDir, filename)
	if err := writeContentFile(content, outputPath); err != nil {
		return "", fmt.Errorf("failed to save summary: %w", err)
	}

	return outputPath, nil
}

// RenameTranscriptFile renames the selected transcript file to a dated format based on the folder date.
// Returns the new filename if renamed, empty string if skipped, or error if failed.
// Skips rename if: already dated, no date in folder path, or transcriptPath not set.
func (p *Processor) RenameTranscriptFile() (string, error) {
	if p.transcriptPath == "" {
		return "", nil
	}

	filename := filepath.Base(p.transcriptPath)

	// Check if already a dated transcript (skip rename)
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-transcript\.txt$`)
	if datePattern.MatchString(filename) {
		return "", nil // Already dated, nothing to do
	}

	// Get date from folder path
	date := p.ExtractDateFromPath()
	if date == "" {
		return "", nil // No date in path, skip rename
	}

	// Build new filename and path
	newFilename := fmt.Sprintf("%s-transcript.txt", date)
	newPath := filepath.Join(p.meetingDir, newFilename)

	// Check if destination already exists
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("cannot rename transcript: %s already exists", newFilename)
	}

	// Rename the file
	if err := os.Rename(p.transcriptPath, newPath); err != nil {
		return "", fmt.Errorf("failed to rename transcript: %w", err)
	}

	// Update the stored path to reflect the new location
	p.transcriptPath = newPath

	return newFilename, nil
}
