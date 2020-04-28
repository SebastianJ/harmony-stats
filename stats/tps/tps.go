package tps

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/SebastianJ/harmony-stats/charts"
	"github.com/SebastianJ/harmony-stats/config"
	sdkRPC "github.com/harmony-one/go-lib/rpc"
)

var (
	targetShards []uint32
)

// BlockResult - statistics for each block
type BlockResult struct {
	ShardID     uint32
	BlockNumber uint64
	TxCount     uint64
	TPS         float64
	Successful  bool
}

// AnalyzeTPS - analyze TPS based on reported txs per every block
func AnalyzeTPS() error {
	if err := setTargetShards(); err != nil {
		return err
	}

	var waitGroup sync.WaitGroup

	for _, shard := range targetShards {
		waitGroup.Add(1)
		go analyzeTPSForShard(shard, &waitGroup)
	}

	waitGroup.Wait()

	return nil
}

func analyzeTPSForShard(shard uint32, parentWaitGroup *sync.WaitGroup) error {
	defer parentWaitGroup.Done()

	node := config.Configuration.Network.API.Shards[shard].Node
	fmt.Printf("Checking tx counts for shard %d\n", shard)

	fromBlockNumber := uint64(0)
	toBlockNumber := uint64(0)
	latestBlockNumber, err := sdkRPC.GetCurrentBlockNumber(node)
	if err != nil {
		return err
	}

	fmt.Printf("latestBlockNumber is now: %d\n", latestBlockNumber)

	if config.TPSArgs.From >= 0 && config.TPSArgs.To >= 0 {
		fromBlockNumber = uint64(config.TPSArgs.From)
		toBlockNumber = uint64(config.TPSArgs.To)
	} else if config.TPSArgs.From >= 0 && config.TPSArgs.To < 0 {
		fromBlockNumber = uint64(config.TPSArgs.From)
		toBlockNumber = latestBlockNumber
	} else if config.TPSArgs.From >= 0 && config.TPSArgs.Count >= 0 {
		fromBlockNumber = uint64(config.TPSArgs.From)
		toBlockNumber = fromBlockNumber + uint64(config.TPSArgs.Count)
	} else if config.TPSArgs.To >= 0 && config.TPSArgs.Count >= 0 {
		toBlockNumber = uint64(config.TPSArgs.To)
		fromBlockNumber = toBlockNumber - uint64(config.TPSArgs.Count)
	} else if config.TPSArgs.From < 0 && config.TPSArgs.To < 0 && config.TPSArgs.Count > 0 {
		toBlockNumber = latestBlockNumber
		fromBlockNumber = toBlockNumber - uint64(config.TPSArgs.Count)
	} else {
		fromBlockNumber = 0
		toBlockNumber = latestBlockNumber
	}

	fmt.Printf("Starting to analyze blocks from block #%d to block #%d for shard %d ...\n", fromBlockNumber, toBlockNumber, shard)

	currentBlockNumber := fromBlockNumber
	totalChecked := 0
	var innerWaitGroup sync.WaitGroup
	blockResults := make(chan BlockResult, toBlockNumber)

	for {
		if currentBlockNumber < toBlockNumber {
			innerWaitGroup.Add(1)

			go blockStatistics(node, shard, currentBlockNumber, blockResults, &innerWaitGroup)
			totalChecked++

			// Wait every <concurrency count> number of blocks before proceeding to queue up more goroutines
			if currentBlockNumber%uint64(config.Configuration.Concurrency) == 0 {
				innerWaitGroup.Wait()
			}

			currentBlockNumber++
		} else {
			break
		}
	}

	innerWaitGroup.Wait()

	close(blockResults)

	results := []BlockResult{}

	for blockResult := range blockResults {
		if blockResult.Successful {
			fmt.Printf("Tx Count for block number %d in shard %d is: %d - TPS is %f\n", blockResult.BlockNumber, blockResult.ShardID, blockResult.TxCount, blockResult.TPS)
			results = append(results, blockResult)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BlockNumber < results[j].BlockNumber
	})

	fileName := fmt.Sprintf("shard-%d-block-%d-to-%d.png", shard, fromBlockNumber, toBlockNumber)
	xAxisData, yAxisData := convertBlockResultsToGraphData(results)

	charts.GenerateGraph(
		fileName,
		"Transactions Per Second",
		"Block #",
		"Transactions Per Second",
		xAxisData,
		yAxisData,
		[]string{
			"Harmony TX/s Report",
			fmt.Sprintf("Network: %s", config.Configuration.Network.Name),
			fmt.Sprintf("Shard: %d", shard),
			fmt.Sprintf("Blocks: %d - %d", fromBlockNumber, toBlockNumber),
		},
	)

	return nil
}

func blockStatistics(node string, shard uint32, blockNumber uint64, blockResults chan<- BlockResult, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	blockResult := BlockResult{}

	fmt.Printf("Checking tx count and tps for block number %d in shard %d (node: %s) ...\n", blockNumber, shard, node)

	txCount, err := sdkRPC.GetTransactionCountByBlockNumber(blockNumber, node)
	if err == nil {
		tps := 0.0
		if txCount > 0 {
			tps = float64(txCount) / float64(config.TPSArgs.BlockTime)
		}

		blockResult.Successful = true
		blockResult.ShardID = shard
		blockResult.BlockNumber = blockNumber
		blockResult.TxCount = txCount
		blockResult.TPS = tps
	} else {
		blockResult.Successful = false
	}

	blockResults <- blockResult
}

func setTargetShards() error {
	shardFlag := strings.ToLower(config.TPSArgs.Shard)

	if shardFlag == "all" {
		for i := uint32(0); i < uint32(config.Configuration.Network.API.ShardCount); i++ {
			targetShards = append(targetShards, i)
		}
	} else {
		shard, err := strconv.Atoi(shardFlag)
		if err != nil {
			return err
		}
		targetShards = append(targetShards, uint32(shard))
	}

	return nil
}

func convertBlockResultsToGraphData(blockResults []BlockResult) (xValues []float64, yValues []float64) {
	xValues = []float64{}
	yValues = []float64{}

	for _, blockResult := range blockResults {
		if blockResult.Successful {
			xValues = append(xValues, float64(blockResult.BlockNumber))
			yValues = append(yValues, blockResult.TPS)
		}
	}

	return xValues, yValues
}
