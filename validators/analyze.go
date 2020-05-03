package validators

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/export"
	"github.com/SebastianJ/harmony-stats/utils"
	sdkDelegation "github.com/harmony-one/go-lib/staking/delegation"
	sdkValidator "github.com/harmony-one/go-lib/staking/validator"
	"github.com/harmony-one/harmony/numeric"
)

// ValidatorResult - wrapper for validator info
type ValidatorResult struct {
	Result  sdkValidator.RPCValidatorResult
	Balance numeric.Dec
	Error   error
}

// Analyze - analyze validators
func Analyze() error {
	fmt.Printf("Looking up validator statistics - network: %s, mode: %s, node: %s\n", config.Configuration.Network.Name, config.Configuration.Network.Mode, config.Configuration.Network.Node)

	rpcValidators, err := sdkValidator.AllInformation(config.Configuration.Network.API.NodeAddress(0), true)
	if err != nil {
		return err
	}

	if config.ValidatorArgs.Elected {
		rpcValidators = applyElectedFilter(rpcValidators)
	}

	validatorResults := []ValidatorResult{}
	for _, rpcValidatorResult := range rpcValidators {
		validatorResults = append(validatorResults, ValidatorResult{Result: rpcValidatorResult})
	}

	if config.ValidatorArgs.Balances {
		validatorResults = lookupValidatorBalances(validatorResults)
	}

	filteredValidatorResults := []ValidatorResult{}
	for _, validatorResult := range validatorResults {
		if applyFilters() {
			if matchesFilter(validatorResult) {
				filteredValidatorResults = append(filteredValidatorResults, validatorResult)
			}
		} else {
			filteredValidatorResults = append(filteredValidatorResults, validatorResult)
		}
	}

	fmt.Printf("Total checked number of validators: %d\n", len(rpcValidators))
	fmt.Printf("Total number of validators matching filter: %d\n", len(filteredValidatorResults))

	switch strings.ToLower(config.Configuration.Export.Format) {
	case "csv":
		csvPath, err := exportToCSV(filteredValidatorResults)
		if err != nil {
			return err
		} else if csvPath != "" {
			fmt.Printf("Successfully exported validator data to %s\n", csvPath)
		}
	//case "json":
	default:
	}

	return nil
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

func lookupValidatorBalances(validatorResults []ValidatorResult) []ValidatorResult {
	validatorsChannel := make(chan ValidatorResult, len(validatorResults))
	var waitGroup sync.WaitGroup

	for index, validatorResult := range validatorResults {
		waitGroup.Add(1)
		go lookupValidatorBalance(validatorResult, validatorsChannel, &waitGroup)

		// Wait every <concurrency count> number of blocks before proceeding to queue up more goroutines
		if index%config.Configuration.Concurrency == 0 {
			waitGroup.Wait()
		}
	}

	waitGroup.Wait()
	close(validatorsChannel)

	validatorResults = []ValidatorResult{}

	for validatorResult := range validatorsChannel {
		validatorResults = append(validatorResults, validatorResult)
	}

	return validatorResults
}

func lookupValidatorBalance(validatorResult ValidatorResult, validatorsChannel chan<- ValidatorResult, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	fmt.Printf("Looking up balance for validator wallet %s\n", validatorResult.Result.Validator.Address)

	totalBalance, err := config.Configuration.Network.API.GetTotalBalance(validatorResult.Result.Validator.Address)
	if err != nil {
		validatorResult.Error = err
		validatorsChannel <- validatorResult
		return
	}

	validatorResult.Balance = totalBalance
	validatorsChannel <- validatorResult
}

func exportToCSV(validatorResults []ValidatorResult) (string, error) {
	fileName := fmt.Sprintf("validators/validators-%s-UTC.csv", utils.FormattedTimeString(time.Now().UTC()))

	rows := [][]string{}

	headers := []string{
		"Name",
		"Address",
		"Identity",
		"BLS Key Count",
		"BLS Keys",
		"Self Delegation",
		"Total Delegation",
		"Lifetime Rewards",
	}

	if config.ValidatorArgs.Balances {
		headers = append(headers, "Wallet Balance")
	}

	rows = append(rows, headers)

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
			}

			if config.ValidatorArgs.Balances && !validatorResult.Balance.IsNil() {
				row = append(row, fmt.Sprintf("%f", validatorResult.Balance))
			}

			rows = append(rows, row)
		}
	}

	csvPath, err := export.ExportCSV(fileName, rows)
	if err != nil {
		return "", err
	}

	return csvPath, nil
}
