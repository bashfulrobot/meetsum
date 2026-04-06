package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Paths struct {
		FileBrowserRootDir string `mapstructure:"file_browser_root_dir"`
		AutomationDir      string `mapstructure:"automation_dir"`
		InstructionsFile   string `mapstructure:"instructions_file"`
	} `mapstructure:"paths"`

	Files struct {
		PovInput string `mapstructure:"pov_input"`
	} `mapstructure:"files"`

	AI struct {
		Command string   `mapstructure:"command"`
		Args    []string `mapstructure:"args"`
	} `mapstructure:"ai"`

	Features struct {
		TraceMode   bool `mapstructure:"trace_mode"`
		FileBrowser bool `mapstructure:"file_browser"`
	} `mapstructure:"features"`

	Logging struct {
		Level  string `mapstructure:"level"`  // debug, info, warn, error
		File   string `mapstructure:"file"`   // path to log file
		Output string `mapstructure:"output"` // screen, file, both
	} `mapstructure:"logging"`

	Skills struct {
		WritingStyle string `mapstructure:"writing_style"`
		Humanizer    string `mapstructure:"humanizer"`
	} `mapstructure:"skills"`

	User struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"user"`
}

var AppConfig *Config

// LoadConfig loads configuration from file
func LoadConfig() error {
	AppConfig = &Config{}

	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")

	// Add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/meetsum")
	viper.AddConfigPath("/etc/meetsum")

	// Set defaults
	homeDir, _ := os.UserHomeDir()
	viper.SetDefault("paths.file_browser_root_dir", filepath.Join(homeDir, "Documents", "Company", "Customers"))
	viper.SetDefault("paths.automation_dir", filepath.Join(homeDir, "Documents", "Company", "automation", "summaries"))
	viper.SetDefault("paths.instructions_file", "Meeting-summary-llm-instructions.md")
	viper.SetDefault("files.pov_input", "pov-input.md")
	viper.SetDefault("skills.writing_style", filepath.Join(homeDir, ".claude", "skills", "writing-style", "writing-style.md"))
	viper.SetDefault("skills.humanizer", filepath.Join(homeDir, ".claude", "skills", "humanizer", "humanizer.md"))
	viper.SetDefault("ai.command", "gemini")
	viper.SetDefault("ai.args", []string{})
	viper.SetDefault("features.trace_mode", false)
	viper.SetDefault("features.file_browser", true)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.file", filepath.Join(homeDir, ".config", "meetsum", "error.log"))
	viper.SetDefault("logging.output", "screen")

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Config file not found, use defaults
	}

	// Unmarshal into struct
	return viper.Unmarshal(AppConfig)
}

// GetInstructionsPath returns the full path to the instructions file
func (c *Config) GetInstructionsPath() string {
	return filepath.Join(c.Paths.AutomationDir, c.Paths.InstructionsFile)
}

// GetPovInputPath returns the full path to the POV input file in a meeting directory
func (c *Config) GetPovInputPath(meetingDir string) string {
	return filepath.Join(meetingDir, c.Files.PovInput)
}

// GetWritingSkillPath returns the path to the best available writing skill file.
// Priority: writing_style > humanizer > "" (none).
func (c *Config) GetWritingSkillPath() string {
	ws := c.expandHome(c.Skills.WritingStyle)
	if ws != "" {
		if _, err := os.Stat(ws); err == nil {
			return ws
		}
	}
	h := c.expandHome(c.Skills.Humanizer)
	if h != "" {
		if _, err := os.Stat(h); err == nil {
			return h
		}
	}
	return ""
}

// expandHome replaces a leading ~/ with the user home directory.
func (c *Config) expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

// GetLogFilePath returns the expanded log file path
func (c *Config) GetLogFilePath() string {
	return c.expandHome(c.Logging.File)
}
