# Release Notes

## v1.0.0 - Initial Release (TBD)

### 🎉 First Official Release

This is the initial release of **meetsum**, a modern Go CLI application that transforms meeting transcripts into structured summaries using AI. Built as a complete rewrite of the original bash script with enhanced features and superior user experience.

### ✨ Features

**Core Functionality:**
- 🤖 **AI-Powered Summaries** - Uses Google Gemini Pro for intelligent meeting summary generation
- 📁 **Smart File Detection** - Automatically finds transcripts and optional context files
- 🎯 **Customer Name Extraction** - Intelligent parsing of directory structure for proper naming
- 📄 **Structured Output** - Generates well-formatted markdown summaries with consistent styling

**User Experience:**
- 🎨 **Beautiful Terminal UI** - Modern interface with forms, styling, and interactive prompts
- 📖 **Markdown Preview** - Built-in preview with rich formatting using glamour
- 🔍 **File Browser** - Interactive directory selection when no path provided
- ⚡ **Progress Indicators** - Clear feedback during AI processing

**Configuration & Setup:**
- ⚙️ **YAML Configuration** - Flexible settings system with sensible defaults
- 🔧 **Dependency Management** - Built-in installation for Homebrew and gemini-cli
- 📋 **System Validation** - Comprehensive checking of requirements and setup
- 🛡️ **Security Warnings** - Clear notifications for any external script execution

**Developer Experience:**
- 🚀 **Single Binary** - No external dependencies except gemini-cli
- 🔄 **Cross-Platform** - Native builds for macOS and Linux (Intel/ARM)
- 📖 **Integrated Help** - Comprehensive documentation and setup guidance
- 🧪 **Robust Testing** - Full CI/CD pipeline with automated releases

### 📦 Installation

#### Pre-built Binaries (Recommended)

Download from the [GitHub Releases page](https://github.com/bashfulrobot/meetsum/releases/latest):

**macOS Apple Silicon:**
```bash
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-darwin-arm64
chmod +x meetsum-darwin-arm64
sudo mv meetsum-darwin-arm64 /usr/local/bin/meetsum
```

**macOS Intel:**
```bash
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-darwin-amd64
chmod +x meetsum-darwin-amd64
sudo mv meetsum-darwin-amd64 /usr/local/bin/meetsum
```

**Linux Intel:**
```bash
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-linux-amd64
chmod +x meetsum-linux-amd64
sudo mv meetsum-linux-amd64 /usr/local/bin/meetsum
```

**Linux ARM:**
```bash
curl -LO https://github.com/bashfulrobot/meetsum/releases/latest/download/meetsum-linux-arm64
chmod +x meetsum-linux-arm64
sudo mv meetsum-linux-arm64 /usr/local/bin/meetsum
```

#### From Source
```bash
git clone https://github.com/bashfulrobot/meetsum.git
cd meetsum
go install
```

### 🚀 Quick Start

1. **Install Dependencies:**
   ```bash
   meetsum install all
   ```

2. **Configure Gemini:**
   ```bash
   meetsum docs gemini  # Opens setup documentation
   ```

3. **Generate Your First Summary:**
   ```bash
   meetsum /path/to/meeting/directory
   ```

### 📁 Directory Structure Requirements

meetsum expects this specific structure:

```
/path/to/Customers/
├── CustomerName/           # Customer folder (used in output filename)
│   └── YYYY-MM-DD/        # Date folder
│       ├── transcript.txt          # Required: Meeting transcript
│       ├── pov-input.md           # Optional: Additional context
│       └── YYYY-MM-DD-CustomerName-cadence-call-summary.md  # Generated
```

### 🛠️ Available Commands

| Command | Description |
|---------|-------------|
| `meetsum [dir]` | Generate meeting summary |
| `meetsum check` | Verify dependencies and configuration |
| `meetsum install all` | Install all required dependencies |
| `meetsum install brew` | Install Homebrew (with security warnings) |
| `meetsum install gemini` | Install Gemini CLI |
| `meetsum docs gemini` | Open Gemini CLI documentation |
| `meetsum docs brew` | Open Homebrew website |
| `meetsum version` | Show version information |

### ⚙️ Configuration

Uses `settings.yaml` configuration files in this search order:
1. `./settings.yaml` (current directory)
2. `~/.config/meetsum/settings.yaml` (user config)
3. `/etc/meetsum/settings.yaml` (system-wide)
4. Custom path via `--config` flag

### 🆚 Improvements Over Original Bash Script

| Feature | Bash Script | meetsum (Go) |
|---------|-------------|--------------|
| **External Dependencies** | 5+ CLI tools (gum, bat, glow, etc.) | 1 tool (gemini-cli) |
| **Installation** | Manual script setup | Single binary + auto-installer |
| **Configuration** | Hard-coded variables | YAML configuration system |
| **Error Handling** | Basic shell errors | Comprehensive error messages |
| **User Interface** | Basic gum prompts | Rich terminal UI with styling |
| **Cross-platform** | Unix/Linux only | macOS, Linux (Intel/ARM) |
| **Validation** | Limited file checking | Complete dependency validation |
| **Documentation** | External files | Integrated help and guidance |

### 🔒 Security Features

- **Clear Warnings** for any external script execution (Homebrew installation)
- **Manual Installation Options** for security-conscious users
- **Transparent Documentation** with links to official sources
- **Minimal External Dependencies** compared to bash version

### 🏗️ Technical Details

**Built with:**
- **Go 1.21+** - Modern, compiled language for performance and reliability
- **Cobra** - Professional CLI framework with rich command structure
- **Viper** - Powerful configuration management with multiple sources
- **Charm Libraries** - Beautiful terminal UI components (huh, lipgloss, glamour, log)
- **bitfield/script** - Go-native shell operations for cross-platform compatibility

**Architecture:**
- **Single Binary Deployment** - No runtime dependencies except gemini-cli
- **Configuration-Driven** - Customizable paths, filenames, and behavior
- **Modular Design** - Clean separation of concerns for maintainability
- **Error Recovery** - Graceful handling of missing dependencies and files

### 📊 Release Artifacts

This release includes:
- ✅ **meetsum-darwin-arm64** - macOS Apple Silicon binary
- ✅ **meetsum-darwin-amd64** - macOS Intel binary
- ✅ **meetsum-linux-amd64** - Linux Intel binary
- ✅ **meetsum-linux-arm64** - Linux ARM binary
- ✅ **checksums.txt** - SHA256 checksums for verification
- ✅ **Source code** - Complete source with build instructions

### 🐛 Known Issues

None at this time. Please report issues on [GitHub Issues](https://github.com/bashfulrobot/meetsum/issues).

### 🙏 Acknowledgments

- Original bash script inspiration and requirements
- [Charm.sh](https://charm.sh) for excellent terminal UI libraries
- [Google Gemini](https://ai.google.dev/) for AI-powered text generation
- [Homebrew](https://brew.sh/) for package management
- Go community for excellent tooling and libraries

---

**Full Changelog**: https://github.com/bashfulrobot/meetsum/commits/v1.0.0
**Download**: https://github.com/bashfulrobot/meetsum/releases/tag/v1.0.0

---

## Future Releases

Future release notes will be automatically generated from git commit messages and will appear here when new versions are tagged and released.

### Release Process

This project uses semantic versioning and automated releases:

1. **Development** - All work happens on feature branches merged to `main`
2. **Release** - Push a semver tag (`git tag v1.0.1 && git push origin v1.0.1`)
3. **Automation** - GitHub Actions builds binaries and creates release
4. **Distribution** - Binaries available immediately on GitHub Releases

See [RELEASING.md](RELEASING.md) for detailed release process documentation.