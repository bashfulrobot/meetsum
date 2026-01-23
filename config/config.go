package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Paths struct {
		FileBrowserRootDir string `mapstructure:"file_browser_root_dir"`
		AutomationDir         string `mapstructure:"automation_dir"`
		InstructionsFile      string `mapstructure:"instructions_file"`
	} `mapstructure:"paths"`

	Files struct {
		Transcript string `mapstructure:"transcript"`
		PovInput   string `mapstructure:"pov_input"`
	} `mapstructure:"files"`

	AI struct {
		Command string `mapstructure:"command"`
	} `mapstructure:"ai"`

	Features struct {
		TraceMode   bool `mapstructure:"trace_mode"`
		FileBrowser bool `mapstructure:"file_browser"`
	} `mapstructure:"features"`

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
	viper.SetDefault("paths.file_browser_root_dir", filepath.Join(os.Getenv("HOME"), "Documents", "Company", "Customers"))
	viper.SetDefault("paths.automation_dir", filepath.Join(os.Getenv("HOME"), "Documents", "Company", "automation", "summaries"))
	viper.SetDefault("paths.instructions_file", "Meeting-summary-llm-instructions.md")
	viper.SetDefault("files.transcript", "transcript.txt")
	viper.SetDefault("files.pov_input", "pov-input.md")
	viper.SetDefault("ai.command", "gemini")
	viper.SetDefault("features.trace_mode", false)
	viper.SetDefault("features.file_browser", true)

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

// GetTranscriptPath returns the full path to the transcript file in a meeting directory
func (c *Config) GetTranscriptPath(meetingDir string) string {
	return filepath.Join(meetingDir, c.Files.Transcript)
}

// GetPovInputPath returns the full path to the POV input file in a meeting directory
func (c *Config) GetPovInputPath(meetingDir string) string {
	return filepath.Join(meetingDir, c.Files.PovInput)
}
