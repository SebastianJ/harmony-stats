package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	sdkNetwork "github.com/harmony-one/go-lib/network"
	sdkNetworkTypes "github.com/harmony-one/go-lib/network/types/network"
	sdkNetworkUtils "github.com/harmony-one/go-lib/network/utils"
	goSdkCommon "github.com/harmony-one/go-sdk/pkg/common"
)

// Configuration - the central configuration
var Configuration Config

// Args is a collection of global/persistent flags parsed using Cobra
var Args PersistentFlags

// TPSArgs is a collection of TPS related flags parsed using Cobra
var TPSArgs TPSFlags

// ValidatorArgs is a collection of validator related flags parsed using Cobra
var ValidatorArgs ValidatorFlags

// Configure - configures the test suite tool using a combination of the YAML config file as well as command arguments
func Configure() (err error) {
	if err := configureNetworkConfig(); err != nil {
		return err
	}

	if err := configureApplicationConfig(); err != nil {
		return err
	}

	ConfigureStylingConfig()

	if err = configureExports(); err != nil {
		return err
	}

	return nil
}

func configureNetworkConfig() (err error) {
	if Args.Network != "" {
		Configuration.Network.Name = Args.Network
	}

	Configuration.Network.Name = sdkNetworkUtils.NormalizedNetworkName(Configuration.Network.Name)
	if Configuration.Network.Name == "" {
		return errors.New("you need to specify a valid network name to use! Valid options: localnet, devnet, testnet, pangaea or mainnet")
	}

	Configuration.Network.Mode = strings.ToLower(Configuration.Network.Mode)
	mode := strings.ToLower(Args.Mode)
	if mode != "" && mode != Configuration.Network.Mode {
		Configuration.Network.Mode = mode
	}

	if len(Args.Nodes) > 0 {
		Configuration.Network.Nodes = Args.Nodes
		Configuration.Network.Node = Configuration.Network.Nodes[0]
	} else {
		Configuration.Network.Nodes = []string{}
		if Args.Node != "" && Args.Node != Configuration.Network.Node {
			Configuration.Network.Node = Args.Node
		} else {
			Configuration.Network.Node = sdkNetworkUtils.ResolveStartingNode(Configuration.Network.Name, Configuration.Network.Mode, 0, Configuration.Network.Nodes)
		}
		Configuration.Network.Nodes = append(Configuration.Network.Nodes, Configuration.Network.Node)
	}

	shards, shardingStructure, err := sdkNetworkTypes.GenerateShardSetup(Configuration.Network.Node, Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Nodes)
	if err != nil {
		return fmt.Errorf("failed to generate network & shard setup for network %s using node %s - error: %s", Configuration.Network.Name, Configuration.Network.Node, err.Error())
	}

	if Configuration.Network.Mode == "api" {
		Configuration.Network.Nodes = []string{}
		for _, shard := range shards {
			Configuration.Network.Nodes = append(Configuration.Network.Nodes, shard.Node)
		}
	}

	Configuration.Network.Timeout = Args.Timeout

	Configuration.Network.API = sdkNetworkTypes.Network{
		Name:              Configuration.Network.Name,
		Mode:              Configuration.Network.Mode,
		Node:              Configuration.Network.Node,
		Shards:            shards,
		ShardingStructure: shardingStructure,
	}

	Configuration.Network.API.Initialize()
	if Configuration.Verbose {
		fmt.Printf("Using network: %s, mode: %s, node: %s\n", Configuration.Network.Name, Configuration.Network.Mode, Configuration.Network.Node)
	}

	return nil
}

func configureApplicationConfig() (err error) {
	Configuration.Concurrency = Args.Concurrency

	Configuration.Verbose = Args.Verbose
	// Set the verbosity level of harmony-sdk
	sdkNetwork.Verbose = Configuration.Verbose

	// Set the verbosity level of go-sdk
	goSdkCommon.DebugRPC = Args.VerboseGoSDK

	return nil
}

func configureExports() error {
	Configuration.Export.Path = filepath.Join(Configuration.BasePath, Args.ExportPath)
	if err := os.MkdirAll(Configuration.Export.Path, 0755); err != nil {
		return err
	}

	if Args.Export != "" {
		Configuration.Export.Format = Args.Export
	}

	return nil
}

// ConfigureStylingConfig - configures the styling and color config
func ConfigureStylingConfig() {
	Configuration.Styling.Header = &color.Style{color.FgLightWhite, color.BgBlack, color.OpBold}
	Configuration.Styling.Info = &color.Style{color.FgLightWhite, color.BgGray, color.OpBold}
	Configuration.Styling.Default = &color.Style{color.OpReset}
	Configuration.Styling.Account = &color.Style{color.FgCyan, color.OpBold}
	Configuration.Styling.Funding = &color.Style{color.FgMagenta, color.OpBold}
	Configuration.Styling.Balance = &color.Style{color.FgLightBlue, color.OpBold}
	Configuration.Styling.Transaction = &color.Style{color.FgYellow, color.OpBold}
	Configuration.Styling.Staking = &color.Style{color.FgLightGreen, color.OpBold}
	Configuration.Styling.Teardown = &color.Style{color.FgGray, color.OpBold}
	Configuration.Styling.Success = &color.Style{color.FgLightWhite, color.BgGreen}
	Configuration.Styling.Warning = &color.Style{color.FgLightWhite, color.BgYellow}
	Configuration.Styling.Error = &color.Style{color.FgLightWhite, color.BgRed}
	Configuration.Styling.Padding = strings.Repeat("\t", 10)
}
