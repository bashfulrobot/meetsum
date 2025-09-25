package deps

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/charmbracelet/huh"
)

const (
	GeminiDocsURL = "https://github.com/google-gemini/gemini-cli?tab=readme-ov-file#-authentication-options"
	HomebrewURL   = "https://brew.sh/"
)

// CheckGeminiInstalled checks if gemini-cli is available
func CheckGeminiInstalled() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

// CheckBrewInstalled checks if Homebrew is available
func CheckBrewInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

// InstallHomebrew installs Homebrew using the official installation script
func InstallHomebrew() error {
	fmt.Println(ui.RenderInfo("🍺 Installing Homebrew..."))
	fmt.Println()
	fmt.Println(ui.RenderWarning("⚠️  Security Notice:"))
	fmt.Println(ui.RenderWarning("This will download and execute a shell script from the internet."))
	fmt.Println(ui.RenderWarning("The script will be run with elevated privileges and may prompt for your password."))
	fmt.Println()
	fmt.Println(ui.RenderInfo("🔗 Alternative: Manual installation"))
	fmt.Printf("Visit %s for manual installation instructions\n", ui.AccentStyle.Render(HomebrewURL))
	fmt.Println()

	var proceed bool
	err := huh.NewConfirm().
		Title("Continue with automatic installation?").
		Description("This will run the official Homebrew installation script").
		Value(&proceed).
		Run()
	if err != nil {
		return err
	}

	if !proceed {
		fmt.Println()
		fmt.Println(ui.RenderInfo("👉 Manual installation:"))
		fmt.Printf("Visit %s to install Homebrew manually\n", ui.AccentStyle.Render(HomebrewURL))
		fmt.Println("Then run 'meetsum install gemini' to continue setup")
		return fmt.Errorf("homebrew installation cancelled - use manual installation")
	}

	// Run the Homebrew installation script
	installScript := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

	fmt.Println(ui.RenderInfo("📥 Downloading and running Homebrew installer..."))
	err = exec.Command("bash", "-c", installScript).Run()
	if err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}

	fmt.Println(ui.RenderSuccess("✅ Homebrew installed successfully!"))
	fmt.Println(ui.RenderInfo("💡 You may need to restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc)"))

	return nil
}

// InstallGemini installs gemini-cli via Homebrew
func InstallGemini() error {
	fmt.Println(ui.RenderInfo("🤖 Installing gemini-cli..."))

	if !CheckBrewInstalled() {
		fmt.Println(ui.RenderError("Homebrew is required to install gemini-cli"))

		var installBrew bool
		err := huh.NewConfirm().
			Title("Would you like to install Homebrew first?").
			Value(&installBrew).
			Run()
		if err != nil {
			return err
		}

		if installBrew {
			if err := InstallHomebrew(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("homebrew is required but not installed")
		}
	}

	var proceed bool
	err := huh.NewConfirm().
		Title("Install gemini-cli via Homebrew?").
		Description("This will run 'brew install gemini-cli'").
		Value(&proceed).
		Run()
	if err != nil {
		return err
	}

	if !proceed {
		return fmt.Errorf("gemini-cli installation cancelled")
	}

	fmt.Println(ui.RenderInfo("📦 Installing gemini-cli via Homebrew..."))

	err = exec.Command("brew", "install", "gemini-cli").Run()
	if err != nil {
		return fmt.Errorf("failed to install gemini-cli: %w", err)
	}

	fmt.Println(ui.RenderSuccess("✅ gemini-cli installed successfully!"))

	// Offer to open documentation
	var openDocs bool
	err = huh.NewConfirm().
		Title("Would you like to open the Gemini CLI documentation?").
		Description("This will help you configure authentication").
		Value(&openDocs).
		Run()
	if err == nil && openDocs {
		OpenGeminiDocs()
	}

	fmt.Println()
	fmt.Println(ui.RenderInfo("🔑 Next steps:"))
	fmt.Println(ui.RenderInfo("1. Configure Gemini authentication (see docs)"))
	fmt.Println(ui.RenderInfo("2. Test with: gemini --help"))
	fmt.Println(ui.RenderInfo("3. Run meetsum again"))

	return nil
}

// OpenGeminiDocs opens the Gemini CLI documentation in the default browser
func OpenGeminiDocs() {
	fmt.Println(ui.RenderInfo("📖 Opening Gemini CLI documentation..."))
	openURL(GeminiDocsURL)
}

// OpenHomebrewSite opens the Homebrew website in the default browser
func OpenHomebrewSite() {
	fmt.Println(ui.RenderInfo("🍺 Opening Homebrew website..."))
	openURL(HomebrewURL)
}

// openURL opens a URL in the default browser
func openURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		fmt.Println(ui.RenderWarning("Cannot auto-open browser on this platform"))
		fmt.Printf("Please visit: %s\n", ui.AccentStyle.Render(url))
		return
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Println(ui.RenderWarning("Failed to open browser automatically"))
		fmt.Printf("Please visit: %s\n", ui.AccentStyle.Render(url))
	}
}

// ValidateGeminiSetup checks if gemini is installed and configured
func ValidateGeminiSetup() error {
	if !CheckGeminiInstalled() {
		return fmt.Errorf("gemini-cli is not installed")
	}

	// Test if gemini is configured by running a simple command
	err := exec.Command("gemini", "--help").Run()
	if err != nil {
		return fmt.Errorf("gemini-cli is installed but may not be properly configured: %w", err)
	}

	return nil
}