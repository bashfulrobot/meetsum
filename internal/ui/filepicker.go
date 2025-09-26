package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bashfulrobot/meetsum/config"
)

type filePickerModel struct {
	filepicker       filepicker.Model
	selectedPath     string
	quitting         bool
	err              error
	transcriptFile   string
}

type clearErrorMsg struct{}

func clearErrorAfter() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m filePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m filePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file or directory?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Check if the selected path is a directory
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			// Selected a directory - check if it contains the transcript file
			transcriptPath := filepath.Join(path, m.transcriptFile)
			if _, err := os.Stat(transcriptPath); err == nil {
				m.selectedPath = path
				m.quitting = true
				return m, tea.Quit
			} else {
				m.err = fmt.Errorf("directory must contain %s", m.transcriptFile)
				return m, clearErrorAfter()
			}
		} else {
			// Selected a file - check the parent directory
			dir := filepath.Dir(path)
			transcriptPath := filepath.Join(dir, m.transcriptFile)
			if _, err := os.Stat(transcriptPath); err == nil {
				m.selectedPath = dir
				m.quitting = true
				return m, tea.Quit
			} else {
				m.err = fmt.Errorf("selected directory must contain %s", m.transcriptFile)
				return m, clearErrorAfter()
			}
		}
	}

	// Did the user select a disabled file?
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = fmt.Errorf("%s is not valid", path)
		return m, clearErrorAfter()
	}

	return m, cmd
}

func (m filePickerModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	s.WriteString(HeaderStyle.Render("📁 Select Meeting Directory"))
	s.WriteString("\n\n")
	s.WriteString(InfoStyle.Render(fmt.Sprintf("Navigate to a directory containing %s", m.transcriptFile)))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("Use ↑/↓ to navigate, Enter to select directory, q to quit"))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		s.WriteString("\n\n")
	}

	s.WriteString(m.filepicker.View())
	s.WriteString("\n")

	return s.String()
}

// SelectDirectory shows a file picker for selecting a meeting directory
func SelectDirectory(startPath string) (string, error) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{} // Allow all files so user can navigate
	fp.CurrentDirectory = startPath
	fp.ShowHidden = false
	fp.DirAllowed = true
	fp.FileAllowed = true
	fp.ShowPermissions = false
	fp.ShowSize = false

	// Get transcript filename from config with fallback
	transcriptFile := "transcript.txt" // default fallback
	if config.AppConfig != nil && config.AppConfig.Files.Transcript != "" {
		transcriptFile = config.AppConfig.Files.Transcript
	}

	m := filePickerModel{
		filepicker:     fp,
		transcriptFile: transcriptFile,
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	// Try different casting approaches for the final model
	switch fm := finalModel.(type) {
	case *filePickerModel:
		if fm.selectedPath != "" {
			return fm.selectedPath, nil
		}
	case filePickerModel:
		if fm.selectedPath != "" {
			return fm.selectedPath, nil
		}
	}

	return "", fmt.Errorf("no directory selected")
}