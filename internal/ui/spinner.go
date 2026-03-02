package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	result  any
	err     error
}

type SpinnerFinishedMsg struct {
	Result any
	Err    error
}

func NewSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = AccentStyle

	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case SpinnerFinishedMsg:
		m.done = true
		m.result = msg.Result
		m.err = msg.Err
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err))
		}
		return SuccessStyle.Render("✅ Complete!")
	}

	return fmt.Sprintf("%s %s", m.spinner.View(), InfoStyle.Render(m.message))
}

// RunWithSpinner executes a function with a spinner display
func RunWithSpinner(message string, fn func() (any, error)) (any, error) {
	model := NewSpinnerModel(message)

	p := tea.NewProgram(model)

	// Run the function in a goroutine
	go func() {
		result, err := fn()
		p.Send(SpinnerFinishedMsg{Result: result, Err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(SpinnerModel); ok {
		return m.result, m.err
	}

	return nil, fmt.Errorf("unexpected model type")
}
