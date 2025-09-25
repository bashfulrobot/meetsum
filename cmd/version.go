package cmd

import (
	"fmt"

	"github.com/bashfulrobot/meetsum/internal/ui"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version, build time, and git commit information for meetsum.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.RenderHeader("🤖 meetsum", "Meeting Summary Generator"))
		fmt.Println()
		fmt.Printf("Version:    %s\n", ui.AccentStyle.Render(version))
		fmt.Printf("Built:      %s\n", ui.SecondaryStyle.Render(buildTime))
		fmt.Printf("Commit:     %s\n", ui.SecondaryStyle.Render(gitCommit))
		fmt.Println()
		fmt.Printf("Repository: %s\n", ui.AccentStyle.Render("https://github.com/bashfulrobot/meetsum"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}