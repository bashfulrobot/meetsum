# 🤖 meetsum

**Meeting Summary Generator CLI Tool**

A modern Go CLI application that transforms meeting transcripts into structured summaries using AI.

[![Go](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Features

- 🎯 **AI-Powered Summaries** - Uses Google Gemini Pro to generate structured meeting summaries
- 📁 **Interactive File Picker** - Visual directory navigation with automatic transcript detection
- 🎨 **Beautiful Terminal UI** - Modern interface with forms, spinners, styling, and markdown preview
- ⚙️ **YAML Configuration** - Flexible configuration system with sensible defaults
- 🔧 **Dependency Management** - Built-in installation and validation of required tools
- 🚀 **Cross-Platform** - Builds for macOS, Linux (Intel/ARM)
- 📖 **Integrated Documentation** - Quick access to setup guides and documentation
- 🔒 **Security-Conscious** - Clear warnings for any external script execution

## ⚠️ Breaking Changes

### Configuration Rename (v0.x.x)

The configuration setting `customers_dir` has been renamed to `file_browser_root_dir` for clarity.

**Migration**: Update your `settings.yaml`:

```yaml
# Before
paths:
  customers_dir: "/path/to/your/Customers"

# After
paths:
  file_browser_root_dir: "/path/to/your/Customers"
```

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
│       ├── <exactly-one>.txt       # Required: exactly one transcript candidate
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
   - Exactly one file with `.txt` extension (case-insensitive) in the meeting directory
   - `Meeting-summary-llm-instructions.md` - Must exist in your automation directory
   - If zero or multiple `.txt` files are present, `meetsum`, the file picker, and `meetsum validate` all fail fast

4. **Optional Files**:
   - `pov-input.md` - Additional context and structure guidance (configurable filename)

### Example Valid Paths

```
✅ /home/user/Documents/Company/Customers/CustomerA/2024-01-15/
✅ /Users/user/Work/Customers/ClientCorp/2024-03-20/
✅ /data/meetings/Customers/BusinessInc/2024-12-01/
```

### Example Invalid Paths

```
❌ /home/user/meetings/CustomerA/2024-01-15/         # Missing "Customers" in path
❌ /home/user/Documents/Customers/2024-01-15/        # Missing customer name folder
❌ /home/user/Documents/Customers/Client Corp/date/  # Spaces in customer name (problematic)
```

### Output Filename Generation

The generated summary filename follows this pattern: `{date}-{customer}-cadence-call-summary.md`

- `{date}` - Current date in YYYY-MM-DD format (when summary is generated)
- `{customer}` - Extracted from the directory path structure
- Example: `2024-01-15-CustomerA-cadence-call-summary.md`

**Note**: The filename format is automatically generated and cannot be customized.

### Path Configuration

You can customize the base paths in your configuration file:

```yaml
paths:
  file_browser_root_dir: "/your/path/to/Customers"
  automation_dir: "/your/path/to/automation/summaries"

files:
  transcript: "transcript.txt"      # Deprecated: transcript source is auto-discovered
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
  file_browser_root_dir: "/path/to/your/Customers"
  automation_dir: "/path/to/automation/summaries"
  instructions_file: "Meeting-summary-llm-instructions.md"

ai:
  command: "gemini"

features:
  trace_mode: false
  file_browser: true
  markdown_preview: true

user:
  name: "Your Name"  # Skip the name prompt; use --ask-name to override
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
| `--ask-name` | Prompt for name even if `user.name` is configured |

## 🏗️ Development

### Prerequisites

- Go 1.25+
- C compiler (gcc, clang, or build-essential)
- [just](https://github.com/casey/just) task runner (optional)

### Building from Source

**Note**: This project requires CGO due to some dependencies. You need a C compiler installed:

```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# macOS (via Xcode Command Line Tools)
xcode-select --install

# Or use clang if available
export CC=clang
```

```bash
git clone https://github.com/bashfulrobot/meetsum.git
cd meetsum
go mod tidy
go build -o meetsum
```

If you encounter CGO errors, ensure you have a C compiler and set the appropriate environment variables:
```bash
# For clang users
CC=clang go build -o meetsum
```

**Note for Contributors**: The project uses GitHub Actions with a matrix build strategy to handle CGO cross-compilation across different platforms (Linux, macOS) and architectures (amd64, arm64). Each platform builds natively to ensure proper CGO support.

### Development Commands (with just)

```bash
# Show available commands
just

# Run deterministic local quality checks (format, vet, tests)
just quality

# Build for current platform
just build

# Run tests
just test

# Build for all platforms
just build-all

# Development workflow
just quality && just build
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

## 🔒 Security Considerations

meetsum takes security seriously:

- **⚠️ Clear Warnings** - All external script execution (Homebrew installation) includes explicit security warnings
- **🔍 Manual Options** - Alternative manual installation paths for security-conscious users
- **📖 Transparent Documentation** - Links to official documentation for all dependencies
- **🛡️ Minimal Dependencies** - Reduced external tool requirements

When installing Homebrew, you'll see warnings like:
```
⚠️  Security Notice:
This will download and execute a shell script from the internet.
The script will be run with elevated privileges and may prompt for your password.

🔗 Alternative: Manual installation
Visit https://brew.sh/ for manual installation instructions
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm.sh](https://charm.sh) for excellent terminal UI libraries
- [Google Gemini](https://ai.google.dev/) for AI-powered text generation
- [Homebrew](https://brew.sh/) for package management
- [Just](https://github.com/casey/just) for task running

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/bashfulrobot/meetsum/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/bashfulrobot/meetsum/discussions)
- 📖 **Gemini Documentation**: Run `meetsum docs gemini` for Gemini CLI setup instructions
