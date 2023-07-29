package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/seungjulee/simple-eth-indexer/internal/server"
	"github.com/seungjulee/simple-eth-indexer/pkg/datastore"
	"github.com/seungjulee/simple-eth-indexer/pkg/logger"
	"github.com/seungjulee/simple-eth-indexer/rpc"
	"github.com/twitchtv/twirp"
	"gopkg.in/yaml.v2"
)

type Config struct {
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

	db, err := datastore.NewSqllite(cfg.SqliteConfig)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
	sv := server.NewExplorerServer(db)
	twirpHandler := rpc.NewSimpleEthExplorerServer(sv, twirp.WithServerHooks(server.NewLoggingServerHooks()))
	logger.Info("starting the server on port 8080")
	if err := http.ListenAndServe(":8080", twirpHandler); err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
}