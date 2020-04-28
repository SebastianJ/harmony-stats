package config

// PersistentFlags represents the persistent flags
type PersistentFlags struct {
	Network      string
	Mode         string
	Node         string
	Nodes        []string
	Timeout      int
	Concurrency  int
	Verbose      bool
	VerboseGoSDK bool
}

// TPSFlags tps related configuration flags
type TPSFlags struct {
	Shard     string
	From      int
	To        int
	Count     int
	BlockTime int
}
