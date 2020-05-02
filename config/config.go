package config

import (
	"github.com/gookit/color"
	sdkAccounts "github.com/harmony-one/go-lib/accounts"
	sdkNetworkTypes "github.com/harmony-one/go-lib/network/types/network"
)

// Config - general config
type Config struct {
	BasePath    string
	Network     Network
	Account     sdkAccounts.Account
	Verbose     bool
	Styling     Styling
	Concurrency int
	Export      Export
}

// Network - represents the network settings group
type Network struct {
	Name    string
	Mode    string
	Node    string
	Nodes   []string
	Shards  int
	API     sdkNetworkTypes.Network
	Timeout int
}

// Export - export settings
type Export struct {
	Path   string
	Format string
}

// Styling - represents settings for styling the log output
type Styling struct {
	Header      *color.Style
	Info        *color.Style
	Default     *color.Style
	Account     *color.Style
	Funding     *color.Style
	Balance     *color.Style
	Transaction *color.Style
	Staking     *color.Style
	Teardown    *color.Style
	Success     *color.Style
	Warning     *color.Style
	Error       *color.Style
	Padding     string
}
