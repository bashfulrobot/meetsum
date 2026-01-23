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
	filepicker     filepicker.Model
	selectedPath   string
	quitting       bool
	err            error
	transcriptFile string
	povInputFile   string
	rootPath       string
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
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "h", "left", "backspace":
			// Block navigation above root
			if m.filepicker.CurrentDirectory == m.rootPath {
				return m, nil
			}
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Ensure we never go above root (safety check)
	if !strings.HasPrefix(m.filepicker.CurrentDirectory, m.rootPath) {
		m.filepicker.CurrentDirectory = m.rootPath
	}

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

	// Current location
	s.WriteString(SecondaryStyle.Render(fmt.Sprintf("Location: %s", m.filepicker.CurrentDirectory)))
	s.WriteString("\n\n")

	// File requirements
	s.WriteString(InfoStyle.Render(fmt.Sprintf("Required: %s", m.transcriptFile)))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render(fmt.Sprintf("Optional: %s", m.povInputFile)))
	s.WriteString("\n\n")

	// Key bindings
	s.WriteString(SecondaryStyle.Render("Navigation:"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  ↑/k      Move up"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  ↓/j      Move down"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  →/l      Enter directory"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  ←/h      Go back"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  Enter    Select directory"))
	s.WriteString("\n")
	s.WriteString(SecondaryStyle.Render("  q/Esc    Quit"))
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
	// Validate startPath exists and is a directory
	if info, err := os.Stat(startPath); err != nil || !info.IsDir() {
		// Return error instead of silently falling back
		return "", fmt.Errorf("file browser root directory does not exist: %s\nPlease update 'file_browser_root_dir' in your settings.yaml", startPath)
	}

	fp := filepicker.New()
	fp.AllowedTypes = []string{} // Allow all files so user can navigate
	fp.CurrentDirectory = startPath
	fp.ShowHidden = false
	fp.DirAllowed = true
	fp.FileAllowed = true
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.AutoHeight = true // Let filepicker auto-size based on terminal

	// Get filenames from config with fallbacks
	transcriptFile := "transcript.txt"
	if config.AppConfig != nil && config.AppConfig.Files.Transcript != "" {
		transcriptFile = config.AppConfig.Files.Transcript
	}
	povInputFile := "pov-input.md"
	if config.AppConfig != nil && config.AppConfig.Files.PovInput != "" {
		povInputFile = config.AppConfig.Files.PovInput
	}

	m := filePickerModel{
		filepicker:     fp,
		transcriptFile: transcriptFile,
		povInputFile:   povInputFile,
		rootPath:       startPath,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	// Get the selected path from the final model
	if fm, ok := finalModel.(filePickerModel); ok {
		if fm.selectedPath != "" {
			return fm.selectedPath, nil
		}
	}

	return "", fmt.Errorf("no directory selected")
}