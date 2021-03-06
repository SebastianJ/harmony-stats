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
	Path         string
	Export       string
	ExportPath   string
}

// TPSFlags tps related configuration flags
type TPSFlags struct {
	Shard     string
	From      int
	To        int
	Count     int
	BlockTime int
}

// ValidatorFlags validator related configuration flags
type ValidatorFlags struct {
	Filter   FilterFlags
	Elected  bool
	Balances bool
}

// FilterFlags - filter validators based on certain criteria
type FilterFlags struct {
	Field string
	Value string
	Mode  string
}
