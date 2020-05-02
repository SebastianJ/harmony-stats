package validators

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/export"
	sdkDelegation "github.com/harmony-one/go-lib/staking/delegation"
	sdkValidator "github.com/harmony-one/go-lib/staking/validator"
	"github.com/harmony-one/harmony/numeric"
)

// ValidatorResult - wrapper for validator info
type ValidatorResult struct {
	Result  sdkValidator.RPCValidatorResult
	Balance numeric.Dec
}

// Analyze - analyze validators
func Analyze() error {
	fmt.Printf("Looking up validator statistics - network: %s, mode: %s, node: %s\n", config.Configuration.Network.Name, config.Configuration.Network.Mode, config.Configuration.Network.Node)

	pageSize := 10
	pages := 1
	var validators []string
	var err error
	message := ""
	rpcClient, err := config.Configuration.Network.API.RPCClient(0)
	if err != nil {
		return err
	}

	if config.ValidatorArgs.Elected {
		message = " elected"
		validators, err = sdkValidator.AllElected(rpcClient)
	} else {
		validators, err = sdkValidator.All(rpcClient)
	}

	if err != nil {
		return err
	}

	validatorCount := len(validators)
	fmt.Println(fmt.Sprintf("Found a total of %d%s validators to look up delegation information for", validatorCount, message))
	pages = calculatePageCount(validatorCount, pageSize)
	totalChecked := 0

	validatorsChannel := make(chan ValidatorResult, validatorCount)
	var waitGroup sync.WaitGroup

	for page := 0; page < pages; page++ {
		for i := 0; i < pageSize; i++ {
			position, ok := processable(page, pageSize, i, validatorCount)
			if ok {
				waitGroup.Add(1)
				address := validators[position]
				go lookupValidator(address, validatorsChannel, &waitGroup)
				totalChecked++
			}
		}

		waitGroup.Wait()
	}

	close(validatorsChannel)

	filteredValidatorResults := []ValidatorResult{}
	for validatorResult := range validatorsChannel {
		if applyFilters() {
			if matchesFilter(validatorResult) {
				filteredValidatorResults = append(filteredValidatorResults, validatorResult)
			}
		} else {
			filteredValidatorResults = append(filteredValidatorResults, validatorResult)
		}
	}

	filterCount := len(filteredValidatorResults)
	fmt.Printf("Total checked number of validators: %d\n", totalChecked)
	fmt.Printf("Total number of validators matching filter: %d\n", filterCount)

	switch strings.ToLower(config.Configuration.Export.Format) {
	case "csv":
		csvPath, err := exportToCSV(filteredValidatorResults)
		if err != nil {
			fmt.Println("Failed to export validator data to CSV")
		} else if csvPath != "" {
			fmt.Printf("Successfully exported validator data to %s\n", csvPath)
		}
	//case "json":
	default:
	}

	return nil
}

func applyFilters() bool {
	return config.ValidatorArgs.Filter.Field != "" && config.ValidatorArgs.Filter.Value != "" && config.ValidatorArgs.Filter.Mode != ""
}

func matchesFilter(validatorResult ValidatorResult) bool {
	currentValue := ""

	switch strings.ToLower(config.ValidatorArgs.Filter.Field) {
	case "website":
		currentValue = validatorResult.Result.Validator.Website
	case "identity":
		currentValue = validatorResult.Result.Validator.Identity
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

func exportToCSV(validatorResults []ValidatorResult) (string, error) {
	rows := [][]string{
		{
			"Name",
			"Address",
			"Identity",
			"BLS Key Count",
			"BLS Keys",
			"Self Delegation",
			"Total Delegation",
			"Lifetime Rewards",
			"Wallet Balance",
		},
	}

	if len(validatorResults) > 0 {
		for _, validatorResult := range validatorResults {
			validator := validatorResult.Result.Validator
			var selfDelegation sdkDelegation.DelegationInfo
			for _, delegation := range validator.Delegations {
				if delegation.DelegatorAddress == validator.Address {
					selfDelegation = delegation
					break
				}
			}

			row := []string{
				validator.Name,
				validator.Address,
				validator.Identity,
				fmt.Sprintf("%d", len(validator.BLSPublicKeys)),
				strings.Join(validator.BLSPublicKeys[:], "\n"),
				fmt.Sprintf("%f", selfDelegation.Amount),
				fmt.Sprintf("%f", validatorResult.Result.TotalDelegation),
				fmt.Sprintf("%f", validatorResult.Result.Lifetime.RewardAccumulated),
				fmt.Sprintf("%f", validatorResult.Balance),
			}

			rows = append(rows, row)
		}
	}

	csvPath, err := export.ExportCSV(rows)
	if err != nil {
		return "", err
	}

	return csvPath, nil
}

func calculatePageCount(totalCount int, pageSize int) int {
	if totalCount > 0 {
		pageNumber := math.RoundToEven(float64(totalCount) / float64(pageSize))
		if math.Mod(float64(totalCount), float64(pageSize)) > 0 {
			return int(pageNumber) + 1
		}

		return int(pageNumber)
	} else {
		return 0
	}
}

func processable(page int, pageSize int, index int, totalCount int) (position int, ok bool) {
	position = ((page * pageSize) + index)
	ok = position <= (totalCount - 1)
	return position, ok
}

func lookupValidator(address string, validatorsChannel chan<- ValidatorResult, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	result, err := sdkValidator.Information(config.Configuration.Network.Node, address)
	if err != nil {
		validatorsChannel <- ValidatorResult{}
		return
	}

	totalBalance, err := config.Configuration.Network.API.GetTotalBalance(address)
	if err != nil {
		validatorsChannel <- ValidatorResult{}
		return
	}

	validatorsChannel <- ValidatorResult{
		Result:  result,
		Balance: totalBalance,
	}
}
