package model

// Transaction is an Ethereum transaction.
type Transaction struct {
	BlockHash string // empty if unconfirmed
	BlockNumber uint64 // empty if unconfirmed

	Hash string `gorm:"primaryKey"`
	Size uint64

	// TxData - Consensus contents of a transaction
	Data []byte
	Gas uint64
	GasPrice string // big int
	Cost uint64
	Nonce uint64
	To string // address
	Value string // big int
}
