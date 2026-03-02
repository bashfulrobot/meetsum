package main

import (
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

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
