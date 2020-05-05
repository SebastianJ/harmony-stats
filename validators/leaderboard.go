package validators

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SebastianJ/harmony-stats/charts"
	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/utils"
	"github.com/wcharczuk/go-chart"
)

// Leaderboard - generate validator leaderboard graph
func Leaderboard() error {
	fmt.Printf("Will generate a graph of the validator leaderboard - network: %s, mode: %s, node: %s\n", config.Configuration.Network.Name, config.Configuration.Network.Mode, config.Configuration.Network.Node)

	validatorResults, err := AcceptableBLS()
	if err != nil {
		return err
	}

	fmt.Printf("Found a total of %d validators eligible to use for the leaderboard\n", len(validatorResults))

	limit := 20

	bars := []chart.Value{}
	fileName := fmt.Sprintf("validators/%s-leaderboard.png", strings.ToLower(config.Configuration.Network.Name))
	for i, validatorResult := range validatorResults {
		if i < limit {
			bar := chart.Value{}
			bar.Label = formatForLabel(validatorResult.Validator.Name)
			rewards := fmt.Sprintf("%f", validatorResult.Lifetime.RewardAccumulated)
			fmt.Printf("Rewards: %s", rewards)
			floatRewards, err := strconv.ParseFloat(rewards, 32)
			if err != nil {
				return err
			}
			fmt.Printf("FloatRewards: %f", floatRewards)

			bar.Value = floatRewards
			bars = append(bars, bar)
		} else {
			break
		}
	}

	if err = charts.GenerateBarChart(fileName, "Open Staking Validator Leaderboard - Lifetime Rewards", bars); err != nil {
		return err
	}

	return nil
}

func formatForLabel(name string) string {
	if len(name) >= 10 {
		name = fixNameWrapping(name)
	}
	name = utils.TruncateString(name, 50)
	return name
}

func fixNameWrapping(name string) string {
	name = strings.ReplaceAll(name, ".", ". ")
	name = strings.ReplaceAll(name, "-", " - ")

	return name
}
