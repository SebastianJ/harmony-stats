package commands

import (
	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/validators"
	"github.com/spf13/cobra"
)

func init() {
	cmdValidators := &cobra.Command{
		Use:   "validators",
		Short: "Validator stats",
		Long:  "Validator statistics etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdAnalyzeValidators := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze validators",
		Long:  "Analyze validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeValidators(cmd)
		},
	}

	config.ValidatorArgs = config.ValidatorFlags{}
	config.ValidatorArgs.Filter = config.Filter{}
	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Field, "filter.field", "", "--filter.field <field>")
	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Value, "filter.value", "", "--filter.value <value>")
	cmdValidators.PersistentFlags().StringVar(&config.ValidatorArgs.Filter.Mode, "filter.mode", "contains", "--filter.mode <mode>")
	cmdValidators.PersistentFlags().BoolVar(&config.ValidatorArgs.Elected, "elected", false, "--elected")
	cmdValidators.AddCommand(cmdAnalyzeValidators)

	RootCmd.AddCommand(cmdValidators)
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
