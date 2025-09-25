package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Color definitions
	HeaderColor    = lipgloss.Color("212")
	InfoColor      = lipgloss.Color("86")
	ErrorColor     = lipgloss.Color("196")
	SuccessColor   = lipgloss.Color("46")
	WarningColor   = lipgloss.Color("220")
	SecondaryColor = lipgloss.Color("240")
	AccentColor    = lipgloss.Color("117")

	// Style definitions
	HeaderStyle = lipgloss.NewStyle().
			Foreground(HeaderColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(HeaderColor).
			Align(lipgloss.Center).
			Width(60).
			Margin(1, 2).
			Padding(1, 2)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	SecondaryStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	AccentStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	InfoBoxStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1, 0)

	FileListStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			Margin(0, 2)
)

// RenderHeader renders the application header
func RenderHeader(title, subtitle string) string {
	content := title
	if subtitle != "" {
		content += "\n" + subtitle
	}
	return HeaderStyle.Render(content)
}

// RenderInfoBox renders an information box
func RenderInfoBox(lines ...string) string {
	content := ""
	for i, line := range lines {
		if i > 0 {
			content += "\n"
		}
		content += line
	}
	return InfoBoxStyle.Render(content)
}

// RenderError renders an error message
func RenderError(message string) string {
	return ErrorStyle.Render("❌ " + message)
}

// RenderSuccess renders a success message
func RenderSuccess(message string) string {
	return SuccessStyle.Render("✅ " + message)
}

// RenderInfo renders an info message
func RenderInfo(message string) string {
	return InfoStyle.Render("ℹ️  " + message)
}

// RenderWarning renders a warning message
func RenderWarning(message string) string {
	return WarningStyle.Render("⚠️  " + message)
}
