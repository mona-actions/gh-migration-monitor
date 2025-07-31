package main

import "github.com/mona-actions/gh-migration-monitor/cmd"

// Version info set by ldflags during build
var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, buildDate)
	cmd.Execute()
}
