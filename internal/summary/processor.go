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

// FindTranscriptFile locates the transcript file in the meeting directory.
// It first looks for the standard transcript.txt, then falls back to
// finding dated transcripts matching YYYY-MM-DD-transcript.txt pattern.
func (p *Processor) FindTranscriptFile() (string, error) {
	// First, try the default transcript.txt
	defaultPath := p.config.GetTranscriptPath(p.meetingDir)
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, nil
	}

	// If not found, look for dated transcript files (YYYY-MM-DD-transcript.txt)
	pattern := filepath.Join(p.meetingDir, "*-transcript.txt")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for transcript files: %w", err)
	}

	// Filter matches to only include valid date-prefixed files
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-transcript\.txt$`)
	for _, match := range matches {
		filename := filepath.Base(match)
		if datePattern.MatchString(filename) {
			// Validate it's a real date
			datePart := filename[:10] // YYYY-MM-DD
			if _, err := time.Parse("2006-01-02", datePart); err == nil {
				return match, nil
			}
		}
	}

	return "", fmt.Errorf("no transcript file found in %s (looked for transcript.txt and YYYY-MM-DD-transcript.txt)", p.meetingDir)
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

// ExtractCustomerName extracts customer name from the meeting directory path
func (p *Processor) ExtractCustomerName() (string, string, error) {
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

	customerNameProper := customerNameRaw
	customerNameUpper := strings.ToUpper(customerNameRaw)

	return customerNameProper, customerNameUpper, nil
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
	customerNameProper, _, err := p.ExtractCustomerName()
	if err != nil {
		return "", err
	}

	date := p.ExtractDateFromPath()
	if date != "" {
		filename := fmt.Sprintf("%s-%s-cadence-call-summary.md", date, customerNameProper)
		return filename, nil
	} else {
		filename := fmt.Sprintf("%s-cadence-call-summary.md", customerNameProper)
		return filename, nil
	}
}

// GenerateSummary processes the meeting and generates the summary
func (p *Processor) GenerateSummary() (string, error) {
	// Load all required content
	instructions, err := p.LoadInstructions()
	if err != nil {
		return "", err
	}

	transcript, err := p.LoadTranscript()
	if err != nil {
		return "", err
	}

	context, err := p.LoadContext()
	if err != nil {
		return "", err
	}

	// Extract date and customer info for the prompt
	customerNameProper, customerNameUpper, err := p.ExtractCustomerName()
	if err != nil {
		return "", err
	}

	date := p.ExtractDateFromPath()
	var titleDate string
	if date != "" {
		titleDate = date
	} else {
		titleDate = "UNDATED"
	}

	// Prepare the prompt
	prompt := fmt.Sprintf(`%s

Process the transcript in transcript.txt and generate a structured meeting summary following the provided instructions. Use the current working directory path to derive the customer name. Write the summary from %s's first-person perspective.

The meeting date should be: %s
The customer name should be: %s (uppercase: %s)

TRANSCRIPT:
%s

%s`, instructions, p.userName, titleDate, customerNameProper, customerNameUpper, transcript, context)

	// Execute AI command with separate stdout/stderr capture
	result, stderr, err := p.executeAICommand(prompt)
	if err != nil {
		p.logCommandError(stderr, err)
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	// Clean the AI output to extract only the markdown content
	cleanedResult := p.cleanAIOutput(result)
	return cleanedResult, nil
}

// executeAICommand runs the AI command and captures stdout/stderr separately
func (p *Processor) executeAICommand(prompt string) (stdout string, stderr string, err error) {
	cmd := exec.Command(p.config.AI.Command)

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
		if strings.Contains(line, "Loaded cached credentials") ||
			strings.Contains(line, "Error executing tool") ||
			strings.Contains(line, "Tool \"write_file\" not found") ||
			strings.Contains(line, "I was unable to create") ||
			strings.Contains(line, "Here is the content") ||
			strings.Contains(line, "You can save it as") {
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

	// Fallback: try to find content after common AI error patterns
	cleanedOutput := output
	errorPatterns := []string{
		"Loaded cached credentials.",
		"Error executing tool write_file:",
		"Tool \"write_file\" not found in registry.",
		"I was unable to create the file directly.",
		"Here is the content for the meeting summary.",
		"You can save it as",
	}

	for _, pattern := range errorPatterns {
		if idx := strings.Index(cleanedOutput, pattern); idx != -1 {
			// Find the end of this line and start from the next line
			nextLine := strings.Index(cleanedOutput[idx:], "\n")
			if nextLine != -1 {
				cleanedOutput = cleanedOutput[idx+nextLine+1:]
			}
		}
	}

	// Remove any remaining markdown code block markers
	cleanedOutput = strings.ReplaceAll(cleanedOutput, "```markdown", "")
	cleanedOutput = strings.ReplaceAll(cleanedOutput, "```", "")

	return strings.TrimSpace(cleanedOutput)
}

// SaveSummary saves the generated summary to a file
func (p *Processor) SaveSummary(content string) (string, error) {
	filename, err := p.GenerateOutputFilename()
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(p.meetingDir, filename)

	// Ensure content ends with a blank line
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	n, err := script.Echo(content).WriteFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to save summary: %w", err)
	}
	_ = n // Ignore bytes written count

	return outputPath, nil
}

// RenameTranscriptFile renames transcript.txt to a dated format based on the folder date.
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
