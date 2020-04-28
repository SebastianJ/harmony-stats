package commands

import (
	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/stats/tps"
	"github.com/spf13/cobra"
)

func init() {
	cmdTps := &cobra.Command{
		Use:   "tps",
		Short: "TPS statistics",
		Long:  "Generate TPS statistics based on transactions per block / block time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeTPS(cmd)
		},
	}

	config.TPSArgs = config.TPSFlags{}
	cmdTps.Flags().StringVar(&config.TPSArgs.Shard, "shard", "all", "--shard <shardID>")
	cmdTps.Flags().IntVar(&config.TPSArgs.From, "from", -1, "--from <blockNumber>")
	cmdTps.Flags().IntVar(&config.TPSArgs.To, "to", -1, "--to <blockNumber>")
	cmdTps.Flags().IntVar(&config.TPSArgs.Count, "count", -1, "--count <count>")
	cmdTps.Flags().IntVar(&config.TPSArgs.BlockTime, "block-time", 8, "--block-time <seconds>")

	RootCmd.AddCommand(cmdTps)
}

func analyzeTPS(cmd *cobra.Command) error {
	if err := config.Configure(); err != nil {
		return err
	}

	if err := tps.AnalyzeTPS(); err != nil {
		return err
	}

	return nil
}
