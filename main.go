package main

import (
	"fmt"
	"os"

	"github.com/bashfulrobot/meetsum/cmd"
)

// Version information - set by build process
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Set version info for use by cobra commands
	cmd.SetVersion(Version, BuildTime, GitCommit)

	// Handle version flag manually since cobra's version handling is limited
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("meetsum version %s\n", Version)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	cmd.Execute()
}
