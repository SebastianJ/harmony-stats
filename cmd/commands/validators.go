package commands

import (
	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/validators"
	"github.com/spf13/cobra"
)

func init() {
	config.ValidatorArgs = config.ValidatorFlags{}
	config.ValidatorArgs.Filter = config.FilterFlags{}

	cmdValidators := &cobra.Command{
		Use:   "validators",
		Short: "Validator stats",
		Long:  "Validator statistics etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Field, "filter.field", "", "--filter.field <field>")
	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Value, "filter.value", "", "--filter.value <value>")
	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Mode, "filter.mode", "contains", "--filter.mode <mode>")
	cmdValidators.PersistentFlags().BoolVar(&config.ValidatorArgs.Elected, "elected", false, "--elected")

	cmdValidators.AddCommand(analyzeCmd())
	cmdValidators.AddCommand(graphsCmd())

	RootCmd.AddCommand(cmdValidators)
}

func analyzeCmd() *cobra.Command {
	cmdAnalyzeValidators := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze validators",
		Long:  "Analyze validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeValidators(cmd)
		},
	}

	cmdAnalyzeValidators.Flags().BoolVar(&config.ValidatorArgs.Balances, "balances", false, "--balances")

	return cmdAnalyzeValidators
}

func analyzeValidators(cmd *cobra.Command) error {
	if err := config.Configure(); err != nil {
		return err
	}

	if err := validators.Analyze(); err != nil {
		return err
	}

	return nil
}

func graphsCmd() *cobra.Command {
	cmdGraphs := &cobra.Command{
		Use:   "graphs",
		Short: "Graph validators",
		Long:  "Graph validator related data",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdDaily := &cobra.Command{
		Use:   "daily",
		Short: "Generate daily graph of validators",
		Long:  "Generate daily graph of validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return graphDailyValidators(cmd)
		},
	}

	cmdLeaderboard := &cobra.Command{
		Use:   "leaderboard",
		Short: "Generate validator leaderboard",
		Long:  "Generate validator leaderboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return graphLeaderboard(cmd)
		},
	}

	cmdGraphs.AddCommand(cmdDaily)
	cmdGraphs.AddCommand(cmdLeaderboard)

	return cmdGraphs
}

func graphDailyValidators(cmd *cobra.Command) error {
	if err := config.Configure(); err != nil {
		return err
	}

	if err := validators.Daily(); err != nil {
		return err
	}

	return nil
}

func graphLeaderboard(cmd *cobra.Command) error {
	if err := config.Configure(); err != nil {
		return err
	}

	if err := validators.Leaderboard(); err != nil {
		return err
	}

	return nil
}
