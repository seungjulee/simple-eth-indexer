package datastore

import (
	"context"

	"github.com/seungjulee/simple-eth-indexer/pkg/datastore/model"
)

type Datastore interface {
	SaveBlock(ctx context.Context, block *model.Block) error
	SaveTXs(ctx context.Context, txs []model.Transaction) error
	SaveEvents(ctx context.Context, events []model.ContractEventLog) error
}