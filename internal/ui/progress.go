package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressModel struct {
	progress progress.Model
	steps    []string
	current  int
	done     bool
	result   interface{}
	err      error
}

type ProgressUpdateMsg struct {
	Step    int
	Message string
}

type ProgressFinishedMsg struct {
	Result interface{}
	Err    error
}

func NewProgressModel(steps []string) ProgressModel {
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 60

	return ProgressModel{
		progress: p,
		steps:    steps,
		current:  0,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case ProgressUpdateMsg:
		m.current = msg.Step
		if msg.Message != "" && msg.Step < len(m.steps) {
			m.steps[msg.Step] = msg.Message
		}
		return m, nil
	case ProgressFinishedMsg:
		m.done = true
		m.result = msg.Result
		m.err = msg.Err
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		model, cmd := m.progress.Update(msg)
		m.progress = model.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m ProgressModel) View() string {
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("%s\n\n%s",
				m.progress.ViewAs(1.0),
				ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
		}
		return fmt.Sprintf("%s\n\n%s",
			m.progress.ViewAs(1.0),
			SuccessStyle.Render("✅ Complete!"))
	}

	// Calculate progress percentage
	percent := float64(m.current) / float64(len(m.steps))

	var output string

	// Show progress bar
	output += m.progress.ViewAs(percent) + "\n\n"

	// Show current step
	if m.current < len(m.steps) {
		output += InfoStyle.Render(fmt.Sprintf("Step %d/%d: %s", m.current+1, len(m.steps), m.steps[m.current]))
	}

	// Show completed steps
	if m.current > 0 {
		output += "\n\n" + SecondaryStyle.Render("Completed:")
		for i := 0; i < m.current && i < len(m.steps); i++ {
			output += fmt.Sprintf("\n  ✓ %s", m.steps[i])
		}
	}

	return output
}

// RunWithProgress executes a function with progress display
func RunWithProgress(steps []string, fn func(updateProgress func(step int, message string)) (interface{}, error)) (interface{}, error) {
	model := NewProgressModel(steps)
	p := tea.NewProgram(model)

	// Run the function in a goroutine
	go func() {
		updateProgress := func(step int, message string) {
			p.Send(ProgressUpdateMsg{Step: step, Message: message})
		}

		result, err := fn(updateProgress)
		p.Send(ProgressFinishedMsg{Result: result, Err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(ProgressModel); ok {
		return m.result, m.err
	}

	return nil, fmt.Errorf("unexpected model type")
}