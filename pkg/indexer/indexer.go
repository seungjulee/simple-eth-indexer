package indexer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/seungjulee/simple-eth-indexer/pkg/datastore"
	"github.com/seungjulee/simple-eth-indexer/pkg/logger"
)

type Indexer interface {
	IndexLastNBlock(context.Context, uint64) error
	IndexBlockAndTXsByNumber(context.Context, uint64) error

	SchedulePeriodicIndex(time.Duration) error
}

type indexer struct {
	ec *ethclient.Client
	db datastore.Datastore
}

func New(ec *ethclient.Client, db datastore.Datastore) Indexer {
	return &indexer{
		ec: ec,
		db: db,
	}
}

func (a *indexer) SchedulePeriodicIndex(interval time.Duration) error {
	logger.Info(fmt.Sprintf("Schedule periodic index every %s", interval))

	blockTicker := time.NewTicker(interval)
	ctx := context.Background()
	for {
		select {
		case <-blockTicker.C:
			a.IndexLastNBlock(ctx, uint64(50))
		}
	}
}

func (a *indexer) IndexLastNBlock(ctx context.Context, n uint64) error {
	if n > 50 {
		errMsg := fmt.Sprintf("you can only index up to last 50 blocks, but got n=%d", n)
		return errors.New(errMsg)
	}
	recentBlockNum, err := a.ec.BlockNumber(ctx)
	if err != nil {
		return err
	}

	for i := recentBlockNum-uint64(n); i <= recentBlockNum; i++ {
		a.IndexBlockAndTXsByNumber(ctx, i)
	}

	return nil
}

func (a *indexer) IndexBlockAndTXsByNumber(ctx context.Context, height uint64) error {
	block, err := a.ec.BlockByNumber(ctx, big.NewInt(int64(height)))
	if err != nil {
		return err
	}

	if err := a.db.SaveBlockAndTXs(ctx, block); err != nil {
		return err
	}

	return nil
}