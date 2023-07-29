package server

import (
	"context"

	"github.com/seungjulee/simple-eth-indexer/pkg/datastore"
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

	return &pb.GetAllEventsByAddressResponse{
	}, nil
}
