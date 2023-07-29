package datastore

import (
	"context"
	"errors"

	"github.com/seungjulee/simple-eth-indexer/pkg/datastore/model"
	"github.com/seungjulee/simple-eth-indexer/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	ormLogger "gorm.io/gorm/logger"
)

type SqliteConfig struct {
	SqlitePath string `yaml:"sqlite_path"`
}

func NewSqllite(cfg *SqliteConfig) (Datastore, error) {
	db, err := gorm.Open(sqlite.Open(cfg.SqlitePath), &gorm.Config{
		Logger: ormLogger.Default.LogMode(ormLogger.Silent),
	  })
	if err != nil {
	  return nil, err
	}
	logger.Info("successfully connected to the db")

	// Migrate the schema
	logger.Info("migrate the schema tables")
	db.AutoMigrate(&model.Block{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.ContractEventLog{})

	return &sqliteDB{
		db: db,
	}, nil
}

type sqliteDB struct {
	db *gorm.DB
}

func (s *sqliteDB) SaveBlock(ctx context.Context, blk *model.Block) error {
	if blk == nil {
		return errors.New("expected non-nil block, but got nil")
	}
	logger.Debug("inserting block into db", zap.Int("block_number", int(blk.Number)), zap.Any("block_hash", blk.Hash))
	if err := s.db.Create(blk).Error; err != nil {
		return err
	}

	return nil
}

func (s *sqliteDB) SaveTXs(ctx context.Context, txs []model.Transaction) error {
	if len(txs) == 0 {
		return errors.New("expected more than 0 txs, but got 0")
	}
	logger.Debug("inserting txs for block into db", zap.Int("block_number", int(txs[0].BlockNumber)), zap.Any("block_hash", txs[0].BlockHash))
	if err := s.db.CreateInBatches(txs, len(txs)).Error; err != nil {
		return err
	}

	return nil
}

func (s *sqliteDB) SaveEvents(ctx context.Context, events []model.ContractEventLog) error {
	if len(events) == 0 {
		return errors.New("expected more than 0 events, but got 0")
	}
	logger.Debug("inserting events for tx for block into db", zap.Int("block_number", int(events[0].BlockNumber)), zap.Any("block_hash", events[0].BlockHash),  zap.Any("tx_hash", events[0].TxHash))
	if err := s.db.CreateInBatches(events, len(events)).Error; err != nil {
		return err
	}

	return nil
}