package datastore

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
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

	return &sqliteDB{
		db: db,
	}, nil
}

type sqliteDB struct {
	db *gorm.DB
}

func (s *sqliteDB) SaveBlockAndTXs(ctx context.Context, block *types.Block) error {
	blk, txs := model.ConvertClientBlockToBlockAndTXs(block)

	logger.Debug("inserting block into db", zap.Any("block_hash", blk.Hash))
	if err := s.db.Create(blk).Error; err != nil {
		return err
	}

	logger.Debug("inserting txs for block into db", zap.Any("block_hash", blk.Hash))
	if err := s.db.CreateInBatches(txs, len(txs)).Error; err != nil {
		return err
	}

	return nil
}

// func (s *sqliteDB) SaveNetInfoAndPeer(netinfo *model.NetInfo, peers []model.Peer) error {
// 	logger.Debug("inserting net_info into db", zap.Any("net_info", netinfo))
// 	if err := s.db.Create(netinfo).Error; err != nil {
// 		return err
// 	}

// 	logger.Debug("inserting peers into db", zap.Any("peers", peers))
// 	if err := s.db.CreateInBatches(peers, len(peers)).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *sqliteDB) GetBlockByHeight(height int64) (*model.Block, error) {
// 	var block model.Block
// 	logger.Debug("get blocks by height", zap.Any("height", height))
// 	if err := s.db.First(&block, "height = ?", height).Error; err != nil {
// 		// fmt.Println(err)
// 		return nil, err
// 	}
// 	return &block, nil
// }
