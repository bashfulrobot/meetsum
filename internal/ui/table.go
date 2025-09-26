package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	table   table.Model
	title   string
	quitting bool
}

func NewTableModel(title string, columns []table.Column, rows []table.Row) TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+2), // +2 for header and padding
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return TableModel{
		table: t,
		title: title,
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	if m.quitting {
		return ""
	}

	var output strings.Builder

	// Title
	if m.title != "" {
		output.WriteString(HeaderStyle.Render(m.title))
		output.WriteString("\n\n")
	}

	// Table
	output.WriteString(m.table.View())
	output.WriteString("\n\n")

	// Help text
	output.WriteString(SecondaryStyle.Render("Use ↑/↓ to navigate, Enter/q to close"))

	return output.String()
}

// ShowTable displays a table and waits for user interaction
func ShowTable(title string, columns []table.Column, rows []table.Row) error {
	model := NewTableModel(title, columns, rows)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}

// FileValidationResult represents the status of a file check
type FileValidationResult struct {
	File        string
	Required    bool
	Found       bool
	Path        string
	Description string
}

// ShowFileValidationTable displays file validation results in a table
func ShowFileValidationTable(results []FileValidationResult) error {
	columns := []table.Column{
		{Title: "File", Width: 25},
		{Title: "Required", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "Path", Width: 40},
		{Title: "Description", Width: 30},
	}

	rows := make([]table.Row, len(results))
	for i, result := range results {
		required := "Optional"
		if result.Required {
			required = "Required"
		}

		status := "❌ Missing"
		if result.Found {
			status = "✅ Found"
		}

		path := result.Path
		if !result.Found {
			path = "Not found"
		}

		rows[i] = table.Row{
			result.File,
			required,
			status,
			path,
			result.Description,
		}
	}

	return ShowTable("📋 File Validation Results", columns, rows)
}

// ConfigItem represents a configuration setting
type ConfigItem struct {
	Category string
	Setting  string
	Value    string
	Default  string
	Description string
}

// ShowConfigurationTable displays configuration settings in a table
func ShowConfigurationTable(config []ConfigItem) error {
	columns := []table.Column{
		{Title: "Category", Width: 15},
		{Title: "Setting", Width: 25},
		{Title: "Current Value", Width: 35},
		{Title: "Default", Width: 25},
		{Title: "Description", Width: 40},
	}

	rows := make([]table.Row, len(config))
	for i, item := range config {
		value := item.Value
		if value == "" {
			value = SecondaryStyle.Render("(empty)")
		}

		defaultVal := item.Default
		if defaultVal == "" {
			defaultVal = SecondaryStyle.Render("(none)")
		}

		rows[i] = table.Row{
			item.Category,
			item.Setting,
			value,
			defaultVal,
			item.Description,
		}
	}

	return ShowTable("⚙️  Configuration Overview", columns, rows)
}