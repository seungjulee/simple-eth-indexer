export GOBIN=$(shell pwd)/bin
export PATH := $(shell pwd)/bin:${PATH}

gen:
	go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

	protoc --proto_path=. --twirp_out=. --go_out=. rpc/*.proto

run_indexer:
	go run ./cmd/indexer/main.go config.yml

run_server:
	go run ./cmd/server/main.go config.yml
