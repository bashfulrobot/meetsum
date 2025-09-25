# 🤖 meetsum

**Meeting Summary Generator CLI Tool**

A modern Go CLI application that transforms meeting transcripts into structured summaries using AI. Built as a complete rewrite of the original bash script with enhanced features, better error handling, and a superior user experience.

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Features

- 🎯 **AI-Powered Summaries** - Uses Google Gemini Pro to generate structured meeting summaries
- 📁 **Intelligent File Detection** - Automatically finds transcripts and context files
- 🎨 **Beautiful Terminal UI** - Modern interface with forms, styling, and markdown preview
- ⚙️ **YAML Configuration** - Flexible configuration system with sensible defaults
- 🔧 **Dependency Management** - Built-in installation and validation of required tools
- 🚀 **Cross-Platform** - Builds for macOS, Linux (Intel/ARM)
- 📖 **Integrated Documentation** - Quick access to setup guides and documentation
- 🔒 **Security-Conscious** - Clear warnings for any external script execution

## 🚀 Quick Start

### 1. Install meetsum

#### Option A: Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [Releases page](https://github.com/bashfulrobot/meetsum/releases/latest):

```bash
# macOS Apple Silicon
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-darwin-arm64
chmod +x meetsum-darwin-arm64
sudo mv meetsum-darwin-arm64 /usr/local/bin/meetsum

# macOS Intel
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-darwin-amd64
chmod +x meetsum-darwin-amd64
sudo mv meetsum-darwin-amd64 /usr/local/bin/meetsum

# Linux Intel
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-linux-amd64
chmod +x meetsum-linux-amd64
sudo mv meetsum-linux-amd64 /usr/local/bin/meetsum

# Linux ARM
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-linux-arm64
chmod +x meetsum-linux-arm64
sudo mv meetsum-linux-arm64 /usr/local/bin/meetsum
```

#### Option B: Install from Source

```bash
git clone https://github.com/bashfulrobot/meetsum.git
cd meetsum
go install
```

#### Verify Installation

```bash
meetsum version
meetsum --help
```

### 2. Install Dependencies

```bash
# Check what's needed
meetsum check

# Install all dependencies automatically
meetsum install all

# Or install individually
meetsum install brew    # Install Homebrew (with security warnings)
meetsum install gemini  # Install Gemini CLI
```

### 3. Configure Gemini

```bash
# Open Gemini documentation for authentication setup
meetsum docs gemini

# Test your setup
meetsum check
```

### 4. Generate Your First Summary

```bash
# Interactive mode - will prompt for meeting directory
meetsum

# Specify directory directly
meetsum /path/to/meeting/directory

# Enable detailed output
meetsum --trace /path/to/meeting/directory
```

## 📁 Directory Structure Requirements

meetsum expects a specific directory structure to properly extract customer names and generate appropriate filenames. **This structure is critical for the application to work correctly.**

### Required Structure

```
/path/to/Customers/
├── CustomerName/           # Customer folder (used in output filename)
│   └── YYYY-MM-DD/        # Date folder (any date format works)
│       ├── transcript.txt          # Required: Meeting transcript
│       ├── pov-input.md           # Optional: Additional context/structure
│       └── YYYY-MM-DD-CustomerName-cadence-call-summary.md  # Generated output
```

### Path Requirements

1. **Customer Name Extraction**: The application extracts the customer name from the directory path structure:
   - Pattern: `/path/to/Customers/{CustomerName}/date/`
   - The folder immediately after "Customers" becomes the customer name in the output filename
   - Customer names should not contain special characters that are invalid in filenames

2. **Date-based Folders**: While not strictly required to be in any specific date format, it's recommended to use:
   - `YYYY-MM-DD` format (e.g., `2024-01-15`)
   - Consistent naming for better organization

3. **Required Files**:
   - `transcript.txt` - The meeting transcript (configurable filename)
   - `Meeting-summary-llm-instructions.md` - Must exist in your automation directory

4. **Optional Files**:
   - `pov-input.md` - Additional context and structure guidance (configurable filename)

### Example Valid Paths

```
✅ /home/user/Documents/Kong/Customers/Acme-Corp/2024-01-15/
✅ /Users/user/Work/Customers/ZillowGroup/2024-03-20/
✅ /data/meetings/Customers/Microsoft/2024-12-01/
```

### Example Invalid Paths

```
❌ /home/user/meetings/Acme-Corp/2024-01-15/        # Missing "Customers" in path
❌ /home/user/Documents/Customers/2024-01-15/       # Missing customer name folder
❌ /home/user/Documents/Customers/Acme Corp/date/   # Spaces in customer name (problematic)
```

### Output Filename Generation

The generated summary filename follows this pattern: `{date}-{customer}-cadence-call-summary.md`

- `{date}` - Current date in YYYY-MM-DD format (when summary is generated)
- `{customer}` - Extracted from the directory path structure
- Example: `2024-01-15-Acme-Corp-cadence-call-summary.md`

### Path Configuration

You can customize the base paths in your configuration file:

```yaml
paths:
  customers_dir: "/your/path/to/Customers/"
  automation_dir: "/your/path/to/automation/summaries"

files:
  transcript: "transcript.txt"      # Customize transcript filename
  pov_input: "pov-input.md"        # Customize context filename
```

## ⚙️ Configuration

meetsum uses YAML configuration files. It searches for settings files in this order:

1. `./settings.yaml` (current directory)
2. `~/.config/meetsum/settings.yaml` (user config)
3. `/etc/meetsum/settings.yaml` (system-wide)
4. Custom path via `--config` flag

### Basic Configuration

```yaml
# settings.yaml
paths:
  customers_dir: "/path/to/your/Customers"
  automation_dir: "/path/to/automation/summaries"
  instructions_file: "Meeting-summary-llm-instructions.md"

ai:
  command: "gemini"

features:
  trace_mode: false
  file_browser: true
  markdown_preview: true
```

### Complete Configuration

See [settings.sample.yaml](settings.sample.yaml) for all available options with detailed comments.

## 🛠️ Commands

### Main Commands

| Command | Description |
|---------|-------------|
| `meetsum [dir]` | Generate meeting summary (interactive if no directory) |
| `meetsum check` | Verify dependencies and configuration |
| `meetsum --help` | Show detailed help and options |

### Installation Commands

| Command | Description |
|---------|-------------|
| `meetsum install all` | Install all required dependencies |
| `meetsum install brew` | Install Homebrew (with security warnings) |
| `meetsum install gemini` | Install Gemini CLI via Homebrew |

### Documentation Commands

| Command | Description |
|---------|-------------|
| `meetsum docs gemini` | Open Gemini CLI documentation |
| `meetsum docs brew` | Open Homebrew website |

### Flags

| Flag | Description |
|------|-------------|
| `--trace` | Enable detailed output, disable loading spinners |
| `--config path` | Use custom configuration file |

## 🏗️ Development

### Prerequisites

- Go 1.21+
- [just](https://github.com/casey/just) task runner (optional)

### Building from Source

```bash
git clone https://github.com/bashfulrobot/meetsum.git
cd meetsum
go mod tidy
go build -o meetsum
```

### Development Commands (with just)

```bash
# Show available commands
just

# Run development checks
just check

# Build for current platform
just build

# Run tests
just test

# Build for all platforms
just build-all

# Development workflow
just check && just test && just build
```

### Project Structure

```
meetsum/
├── cmd/                 # CLI commands and subcommands
│   ├── root.go         # Main command logic
│   ├── check.go        # Dependency validation
│   ├── install.go      # Installation commands
│   └── docs.go         # Documentation commands
├── config/             # Configuration management
│   └── config.go       # YAML config handling
├── internal/
│   ├── deps/           # Dependency checking and installation
│   ├── summary/        # Core summary processing logic
│   └── ui/            # Terminal UI styling and components
├── config.yaml         # Default configuration
├── config.sample.yaml  # Commented configuration example
├── justfile           # Task runner configuration
└── main.go            # Application entry point
```

## 🔧 Dependencies

### External Dependencies

- **gemini-cli** - Google Gemini CLI tool for AI text generation
  - Installation: `brew install gemini-cli`
  - Documentation: https://github.com/google-gemini/gemini-cli

### Internal Dependencies (Go modules)

- `spf13/cobra` - CLI framework
- `spf13/viper` - Configuration management
- `bitfield/script` - Shell operations in Go
- `charmbracelet/huh` - Terminal forms and input
- `charmbracelet/lipgloss` - Terminal styling
- `charmbracelet/glamour` - Markdown rendering
- `charmbracelet/log` - Structured logging

## 🆚 Comparison with Original Bash Script

| Feature | Bash Script | Go CLI (meetsum) |
|---------|-------------|------------------|
| **Dependencies** | 5+ external tools (gum, bat, glow, etc.) | 1 external tool (gemini) |
| **Installation** | Manual script setup | Single binary + auto-install |
| **Configuration** | Hard-coded variables | YAML config files |
| **Error Handling** | Basic shell error handling | Comprehensive error messages |
| **User Interface** | Basic gum prompts | Rich terminal UI with styling |
| **Cross-platform** | Unix/Linux only | macOS, Linux (Intel/ARM) |
| **Validation** | Limited file checking | Complete dependency validation |
| **Documentation** | Separate files | Integrated help and docs |

## 🔒 Security Considerations

meetsum takes security seriously:

- **⚠️ Clear Warnings** - All external script execution (Homebrew installation) includes explicit security warnings
- **🔍 Manual Options** - Alternative manual installation paths for security-conscious users
- **📖 Transparent Documentation** - Links to official documentation for all dependencies
- **🛡️ Minimal Dependencies** - Reduced external tool requirements compared to bash version

When installing Homebrew, you'll see warnings like:
```
⚠️  Security Notice:
This will download and execute a shell script from the internet.
The script will be run with elevated privileges and may prompt for your password.

🔗 Alternative: Manual installation
Visit https://brew.sh/ for manual installation instructions
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and checks (`just check && just test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation for user-facing changes
- Use `just check` before committing to ensure code quality

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original bash script inspiration and requirements
- [Charm.sh](https://charm.sh) for excellent terminal UI libraries
- [Google Gemini](https://ai.google.dev/) for AI-powered text generation
- [Homebrew](https://brew.sh/) for package management
- [Just](https://github.com/casey/just) for task running

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/bashfulrobot/meetsum/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/bashfulrobot/meetsum/discussions)
- 📖 **Documentation**: Run `meetsum docs gemini` or visit the repository

---

**Built with ❤️ and Go**

Made to transform your meeting transcripts into actionable insights with the power of AI and a delightful CLI experience.