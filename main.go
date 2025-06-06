package main

import "esh-cli/cmd"

var version = "dev" // This will be set by ldflags during build

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
