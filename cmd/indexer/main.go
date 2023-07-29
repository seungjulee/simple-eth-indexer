package main

import (
	"context"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/seungjulee/simple-eth-indexer/pkg/datastore"
	"github.com/seungjulee/simple-eth-indexer/pkg/indexer"
	"github.com/seungjulee/simple-eth-indexer/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RPCEndpoint string `yaml:"rpc_endpoint"`
	SqliteConfig *datastore.SqliteConfig `yaml:"sqlite_config"`
}

func ReadConfig(path string) *Config {
    var config Config

    // Open YAML file
    file, err := os.Open(path)
    if err != nil {
        logger.Fatal(err.Error())
    }
    defer file.Close()

    // Decode YAML file to struct
    if file != nil {
        decoder := yaml.NewDecoder(file)
        if err := decoder.Decode(&config); err != nil {
            logger.Fatal(err.Error())
        }
    }

    return &config
}

func main() {
	if len(os.Args) != 2 {
		err := errors.New("add the path for config.yaml. ex) go run cmd/indexer/main.go ./config.yml")
		logger.Fatal(err.Error())
		panic(err)
	}
	cfg := ReadConfig(os.Args[1])

	logger.Info("initializing app with config", zap.Any("config", cfg))
	rpcClient, err := rpc.Dial(cfg.RPCEndpoint)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	defer ethClient.Close()

	db, err := datastore.NewSqllite(cfg.SqliteConfig)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}

	ctx := context.Background()
	a := indexer.New(ethClient, db)
	err = a.IndexLastNBlock(ctx, 50)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
}