package validators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SebastianJ/harmony-stats/config"
	sdkValidator "github.com/harmony-one/go-lib/staking/validator"
)

// All - return all validators
func All() (validatorResults []sdkValidator.RPCValidatorResult, err error) {
	fmt.Printf("Looking up validators - network: %s, mode: %s, node: %s\n", config.Configuration.Network.Name, config.Configuration.Network.Mode, config.Configuration.Network.Node)

	validatorResults, err = sdkValidator.AllInformation(config.Configuration.Network.API.NodeAddress(0), true)
	if err != nil {
		return validatorResults, err
	}

	sort.Slice(validatorResults, func(i, j int) bool {
		return validatorResults[i].Lifetime.RewardAccumulated.GT(validatorResults[j].Lifetime.RewardAccumulated)
	})

	return validatorResults, nil
}

// Elected - return all elected validators
func Elected() (validatorResults []sdkValidator.RPCValidatorResult, err error) {
	validatorResults, err = All()
	if err != nil {
		return validatorResults, err
	}

	validatorResults = applyElectedFilter(validatorResults)

	return validatorResults, nil
}

// AcceptableBLS - return all validators with the correct BLS key limit (1)
func AcceptableBLS() (validatorResults []sdkValidator.RPCValidatorResult, err error) {
	validatorResults, err = All()
	if err != nil {
		return validatorResults, err
	}

	validatorResults = applyAllowedBLSFilter(validatorResults)

	return validatorResults, nil
}

// Filtered - return all validators filtered by certain criteria
func Filtered() (validatorResults []sdkValidator.RPCValidatorResult, err error) {
	validatorResults, err = All()
	if err != nil {
		return validatorResults, err
	}

	if config.ValidatorArgs.Elected {
		validatorResults = applyElectedFilter(validatorResults)
	}

	if applyFilters() {
		filteredValidatorResults := []sdkValidator.RPCValidatorResult{}
		for _, validatorResult := range validatorResults {
			if matchesFilter(validatorResult) {
				filteredValidatorResults = append(filteredValidatorResults, validatorResult)
			}
		}

		validatorResults = filteredValidatorResults
	}

	return validatorResults, nil
}

func applyElectedFilter(validatorResults []sdkValidator.RPCValidatorResult) []sdkValidator.RPCValidatorResult {
	electedValidators := []sdkValidator.RPCValidatorResult{}

	for _, validatorResult := range validatorResults {
		if validatorResult.CurrentlyInCommittee {
			electedValidators = append(electedValidators, validatorResult)
		}
	}

	return electedValidators
}

func applyAllowedBLSFilter(validatorResults []sdkValidator.RPCValidatorResult) []sdkValidator.RPCValidatorResult {
	allowedBLSValidators := []sdkValidator.RPCValidatorResult{}

	for _, validatorResult := range validatorResults {
		if len(validatorResult.Validator.BLSPublicKeys) == 1 {
			allowedBLSValidators = append(allowedBLSValidators, validatorResult)
		}
	}

	return allowedBLSValidators
}

func applyFilters() bool {
	return config.ValidatorArgs.Filter.Field != "" && config.ValidatorArgs.Filter.Value != "" && config.ValidatorArgs.Filter.Mode != ""
}

func matchesFilter(validatorResult sdkValidator.RPCValidatorResult) bool {
	currentValue := ""

	switch strings.ToLower(config.ValidatorArgs.Filter.Field) {
	case "website":
		currentValue = validatorResult.Validator.Website
	case "identity":
		currentValue = validatorResult.Validator.Identity
	default:
		currentValue = ""
	}

	currentValue = strings.ToLower(currentValue)
	expectedValue := strings.ToLower(config.ValidatorArgs.Filter.Value)

	if currentValue != "" {
		if config.ValidatorArgs.Filter.Mode == "equals" {
			return currentValue == expectedValue
		} else if config.ValidatorArgs.Filter.Mode == "contains" {
			return strings.Contains(currentValue, expectedValue)
		}
	}

	return true
}
