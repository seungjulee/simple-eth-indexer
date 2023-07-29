package model

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	Hash string `gorm:"primaryKey"`
	Size uint64

	// Header
	ParentHash string `gorm:"uniqueIndex"`
	UncleHash string
	Coinbase string
	Root string
	TxHash string
	ReceiptHash string
	Difficulty string // big int
	Number uint64 `gorm:"index,sort:desc"` // big int
	GasLimit uint64
	GasUsed uint64
	Time time.Time
	Extra string
	Nonce uint64
	BaseFee string // big int
}

func ConvertClientBlockToBlockAndTXs(b *types.Block) (*Block, []Transaction) {
	block := &Block{
		Hash: b.Header().Hash().Hex(),
		Size: b.Size(),
		ParentHash: b.Header().ParentHash.Hex(),
		UncleHash: b.UncleHash().Hex(),
		Coinbase: b.Coinbase().Hex(),
		Root: b.Root().Hex(),
		TxHash: b.TxHash().Hex(),
		ReceiptHash: b.ReceiptHash().Hex(),
		Difficulty: b.Difficulty().String(),
		Number: b.Number().Uint64(),
		GasLimit: b.GasLimit(),
		GasUsed: b.GasUsed(),
		Time: time.Unix(int64(b.Time()), 0),
		Extra: string(b.Extra()),
		Nonce: b.Header().Nonce.Uint64(),
		BaseFee: b.BaseFee().String(),
	}
	transactions := []Transaction{}
	for _, tx := range b.Transactions() {
		t := Transaction{
			BlockHash: b.Header().Hash().Hex(),
			BlockNumber: b.NumberU64(),
			Hash: tx.Hash().Hex(),
			Size: tx.Size(),
			Data: tx.Data(),
			Gas: tx.Gas(),
			GasPrice: tx.GasPrice().String(),
			Cost: tx.Cost().Uint64(),
			Nonce: tx.Nonce(),
			Value: tx.Value().String(),
		}
		// For contract-creation transactions, To returns nil.
		if tx.To() != nil {
			t.To = tx.To().Hex()
		}
		transactions = append(transactions, t)
	}
	return block, transactions
}
