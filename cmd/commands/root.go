package commands

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/SebastianJ/harmony-stats/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// VersionWrap - version displayed in case of errors
	VersionWrap = fmt.Sprintf("%s/%s-%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// RootCmd - main entry point for Cobra commands
	RootCmd = &cobra.Command{
		Use:          "stats",
		Short:        "Harmony stats",
		SilenceUsage: true,
		Long:         "Harmony stats - generate stats and graphs for Harmony networks",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
)

func init() {
	config.Args = config.PersistentFlags{}
	RootCmd.PersistentFlags().StringVar(&config.Args.Network, "network", "stressnet", "--network <name>")
	RootCmd.PersistentFlags().StringVar(&config.Args.Mode, "mode", "api", "--mode <mode>")
	RootCmd.PersistentFlags().StringVar(&config.Args.Node, "node", "", "--node <node>")
	RootCmd.PersistentFlags().StringSliceVar(&config.Args.Nodes, "nodes", []string{}, "--nodes node1,node2")
	RootCmd.PersistentFlags().IntVar(&config.Args.Timeout, "timeout", 60, "--timeout <timeout>")
	RootCmd.PersistentFlags().IntVar(&config.Args.Concurrency, "concurrency", 100, "<concurrency>")
	RootCmd.PersistentFlags().BoolVar(&config.Args.Verbose, "verbose", false, "--verbose")
	RootCmd.PersistentFlags().BoolVar(&config.Args.VerboseGoSDK, "verbose-go-sdk", false, "--verbose-go-sdk")

	RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "SebastianJ (C) 2020. %v, version %s/%s-%s\n", path.Base(os.Args[0]), runtime.Version(), runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
			return nil
		},
	})
}

// ParseArgs - parse arguments using Cobra
func ParseArgs() {
	RootCmd.SilenceErrors = true
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(errors.Wrapf(err, "commit: %s, error", VersionWrap).Error())
		os.Exit(1)
	}
}
