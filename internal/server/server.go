package server

import (
	"context"
	"strings"

	"github.com/seungjulee/simple-eth-indexer/pkg/datastore"
	"github.com/seungjulee/simple-eth-indexer/pkg/datastore/model"
	pb "github.com/seungjulee/simple-eth-indexer/rpc"
	"github.com/twitchtv/twirp"
)

// Server implements SimpleEthExplorer rpc service
type Server struct {
	ds datastore.Datastore
}

func NewExplorerServer(ds datastore.Datastore) pb.SimpleEthExplorer {
	return &Server{
		ds: ds,
	}
}

func (s *Server) GetAllEventsByAddress(ctx context.Context, req *pb.GetAllEventsByAddressRequest) (*pb.GetAllEventsByAddressResponse, error) {
    if req.Address == "" {
        return nil, twirp.InvalidArgumentError("address", "'address' is empty")
    }

	events, startBlock, endBlock, err := s.ds.GetAllEventsByAddress(ctx, req.Address)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	respEvents := []*pb.Event{}
	for _, e := range events {
		respEvents = append(respEvents, ConvertModelEventToProtoEvent(&e))
	}

	return &pb.GetAllEventsByAddressResponse{
		StartBlock: int64(startBlock),
		EndBlock: int64(endBlock),
		Events: respEvents,
	}, nil
}

func ConvertModelEventToProtoEvent(evt *model.ContractEventLog) *pb.Event {
	event := &pb.Event{
		Address: evt.Address,
		Topics: strings.Split(evt.Topics, ","),
		// TODO: proto: field rpc.Event.data contains invalid UTF-8
		// Data: string(evt.Data),
		BlockNumber: evt.BlockNumber,
		TxHash: evt.TxHash,
		TxIndex: uint64(evt.TxIndex),
		BlockHash: evt.BlockHash,
		Index: uint64(evt.Index),
		Removed: evt.Removed,
	}

	return event
}