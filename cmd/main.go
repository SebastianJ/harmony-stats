package main

import (
	"log"
	"os"
	"path/filepath"

	cmd "github.com/SebastianJ/harmony-stats/cmd/commands"
	"github.com/SebastianJ/harmony-stats/config"
)

func main() {
	// Force usage of Go's own DNS implementation
	os.Setenv("GODEBUG", "netdns=go")

	if err := execute(); err != nil {
		log.Fatalln(err)
	}
}

func execute() error {
	cmd.ParseArgs()

	basePath, err := filepath.Abs(config.Args.Path)
	if err != nil {
		return err
	}

	config.Configuration.BasePath = basePath

	return nil
}
