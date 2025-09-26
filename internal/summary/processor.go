package summary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bashfulrobot/meetsum/config"
	"github.com/bitfield/script"
	"github.com/charmbracelet/log"
)

type Processor struct {
	config     *config.Config
	logger     *log.Logger
	userName   string
	meetingDir string
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

// ValidateRequiredFiles checks if all required files exist
func (p *Processor) ValidateRequiredFiles() error {
	// Check transcript file
	transcriptPath := p.config.GetTranscriptPath(p.meetingDir)
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		return fmt.Errorf("transcript.txt not found in %s", p.meetingDir)
	}

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
	transcriptPath := p.config.GetTranscriptPath(p.meetingDir)
	content, err := script.File(transcriptPath).String()
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

// GenerateOutputFilename creates the output filename
func (p *Processor) GenerateOutputFilename() (string, error) {
	customerNameProper, _, err := p.ExtractCustomerName()
	if err != nil {
		return "", err
	}

	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s-cadence-call-summary.md", date, customerNameProper)

	return filename, nil
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

	// Prepare the prompt
	prompt := fmt.Sprintf(`%s

Process the transcript in transcript.txt and generate a structured meeting summary following the provided instructions. Use the current working directory path to derive the customer name. Write the summary from %s's first-person perspective.

TRANSCRIPT:
%s

%s`, instructions, p.userName, transcript, context)

	// Execute AI command
	result, err := script.Echo(prompt).Exec(p.config.AI.Command).String()
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	// Clean the AI output to extract only the markdown content
	cleanedResult := p.cleanAIOutput(result)
	return cleanedResult, nil
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
		if !foundMarkdownStart && strings.Contains(line, "_SUMMARY_") || strings.HasPrefix(line, "*_") {
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
