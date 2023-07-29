package datastore

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type Datastore interface {
	SaveBlockAndTXs(ctx context.Context, block *types.Block) error
}