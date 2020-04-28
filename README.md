# Harmony Stats
A tool for generating statistics and graphs for Harmony networks

## Installation

```
rm -rf harmony-stats && mkdir -p harmony-stats && cd harmony-stats
bash <(curl -s -S -L https://raw.githubusercontent.com/SebastianJ/harmony-stats/master/scripts/install.sh)
```

## Usage

### Generate Transactions Per Second reports

Generate a TPS graph based on the last x blocks:
```
./stats tps --network NETWORK --shard SHARD_ID --count COUNT
```

```
$ ./stats tps --help
Generate TPS statistics based on transactions per block / block time

Usage:
  stats tps [flags]

Flags:
      --block-time int   --block-time <seconds> (default 8)
      --count int        --count <count> (default -1)
      --from int         --from <blockNumber> (default -1)
  -h, --help             help for tps
      --shard string     --shard <shardID> (default "all")
      --to int           --to <blockNumber> (default -1)

Global Flags:
      --concurrency int   <concurrency> (default 100)
      --mode string       --mode <mode> (default "api")
      --network string    --network <name> (default "stressnet")
      --node string       --node <node>
      --nodes strings     --nodes node1,node2
      --timeout int       --timeout <timeout> (default 60)
      --verbose           --verbose
      --verbose-go-sdk    --verbose-go-sdk
```
