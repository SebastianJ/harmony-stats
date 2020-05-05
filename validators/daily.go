package validators

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/SebastianJ/harmony-stats/charts"
	"github.com/SebastianJ/harmony-stats/config"
	"github.com/elliotchance/orderedmap"
	sdkRPC "github.com/harmony-one/go-lib/rpc"
	sdkValidator "github.com/harmony-one/go-lib/staking/validator"
)

var (
	timeFormat string = "2006-01-02"
)

// Daily - generate validator graph
func Daily() error {
	fmt.Printf("Will generate a graph over daily validators - network: %s, mode: %s, node: %s\n", config.Configuration.Network.Name, config.Configuration.Network.Mode, config.Configuration.Network.Node)

	validatorResults, err := All()
	if err != nil {
		return err
	}

	blockNumbers := []uint64{}
	blockNumberValidatorCountMapping := identifyValidatorCountPerBlock(validatorResults)
	for el := blockNumberValidatorCountMapping.Front(); el != nil; el = el.Next() {
		blockNumber := el.Key.(uint64)
		validatorCount := el.Value.(int)
		fmt.Printf("BlockNumber %d - number of created validators: %d\n", blockNumber, validatorCount)
		blockNumbers = append(blockNumbers, blockNumber)
	}

	fmt.Printf("Retrieving block information for %d block(s)\n", len(blockNumbers))

	blocks := retrieveBlocks(config.Configuration.Network.API.NodeAddress(0), blockNumbers)
	totalCount := 0
	xAxisData := []time.Time{}
	yAxisData := []float64{}
	validatorCountPerDate := identifyValidatorCountPerDate(blocks, blockNumberValidatorCountMapping)

	for el := validatorCountPerDate.Front(); el != nil; el = el.Next() {
		dateString := el.Key.(string)
		validatorCount := el.Value.(int)

		date, err := time.Parse(timeFormat, dateString)
		if err != nil {
			return err
		}

		xAxisData = append(xAxisData, date)
		yAxisData = append(yAxisData, float64(validatorCount))

		totalCount = totalCount + validatorCount
		fmt.Printf("Date %s - number of created validators: %d\n", date, validatorCount)
	}

	fmt.Printf("Total number of created validators: %d\n", totalCount)

	fileName := fmt.Sprintf("validators/%s-daily.png", strings.ToLower(config.Configuration.Network.Name))
	err = charts.GenerateTimeSeriesChart(
		fileName,
		"Validator Count",
		"Date",
		"",
		xAxisData,
		yAxisData,
		[]string{
			fmt.Sprintf("Validators: %d total", totalCount),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func identifyValidatorCountPerBlock(validatorResults []sdkValidator.RPCValidatorResult) *orderedmap.OrderedMap {
	blockNumberValidatorCountMapping := orderedmap.NewOrderedMap()
	for _, validatorResult := range validatorResults {
		if validatorResult.Validator.CreationHeight >= 0 {
			blockNumber := uint64(validatorResult.Validator.CreationHeight)
			value, exists := blockNumberValidatorCountMapping.Get(blockNumber)

			count := 1
			if exists {
				count = value.(int)
				count++
			}

			blockNumberValidatorCountMapping.Set(blockNumber, count)
		}
	}

	return blockNumberValidatorCountMapping
}

func identifyValidatorCountPerDate(blocks []sdkRPC.BlockInfo, blockNumberValidatorCountMapping *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	dateCounts := orderedmap.NewOrderedMap()
	for _, block := range blocks {
		if block.BlockNumber > 0 {
			date := block.Timestamp.Format(timeFormat)

			rawBlockValidatorCount, blockValExists := blockNumberValidatorCountMapping.Get(block.BlockNumber)
			blockValidatorCount := 0
			if blockValExists {
				blockValidatorCount = rawBlockValidatorCount.(int)
			}

			fmt.Printf("Block number: %d, date: %s, validator count: %d, block.Timestamp: %+v\n", block.BlockNumber, date, blockValidatorCount, block.Timestamp)

			value, exists := dateCounts.Get(date)
			validatorCount := blockValidatorCount
			if exists {
				validatorCount = value.(int)
				validatorCount = validatorCount + blockValidatorCount
			}

			dateCounts.Set(date, validatorCount)
		}
	}

	return dateCounts
}

func retrieveBlocks(node string, blockNumbers []uint64) (blockResults []sdkRPC.BlockInfo) {
	blocksChannel := make(chan sdkRPC.BlockInfo, len(blockNumbers))
	var waitGroup sync.WaitGroup

	for index, blockNumber := range blockNumbers {
		waitGroup.Add(1)
		go lookupBlockInfo(node, blockNumber, blocksChannel, &waitGroup)

		// Wait every <concurrency count> number of blocks before proceeding to queue up more goroutines
		if index%config.Configuration.Concurrency == 0 {
			waitGroup.Wait()
		}
	}

	waitGroup.Wait()
	close(blocksChannel)

	blockResults = []sdkRPC.BlockInfo{}
	for blockResult := range blocksChannel {
		blockResults = append(blockResults, blockResult)
	}

	sort.Slice(blockResults, func(i, j int) bool {
		return blockResults[i].BlockNumber < blockResults[j].BlockNumber
	})

	return blockResults
}

func lookupBlockInfo(node string, blockNumber uint64, blocksChannel chan<- sdkRPC.BlockInfo, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	fmt.Printf("Looking up block information for block %d\n", blockNumber)

	blockInfo, err := sdkRPC.GetBlockByNumber(blockNumber, false, node)
	if err != nil {
		blocksChannel <- sdkRPC.BlockInfo{}
		return
	}

	blocksChannel <- blockInfo
}
