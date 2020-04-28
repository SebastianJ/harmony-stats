package main

import (
	"os"

	cmd "github.com/SebastianJ/harmony-stats/cmd/commands"
)

func main() {
	// Force usage of Go's own DNS implementation
	os.Setenv("GODEBUG", "netdns=go")

	cmd.ParseArgs()
}
