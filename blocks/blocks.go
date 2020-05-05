package blocks

// BlockResult - statistics for each block
type BlockResult struct {
	ShardID     uint32
	BlockNumber uint64
	TxCount     uint64
	TPS         float64
	Successful  bool
}
