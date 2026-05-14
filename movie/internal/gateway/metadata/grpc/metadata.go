package grpc

import (
	"context"

	"movieexamplekhubaib.com/gen"
	"movieexamplekhubaib.com/internal/grpcutil"
	"movieexamplekhubaib.com/metadata/pkg/model"
	"movieexamplekhubaib.com/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Return movie metadata by movie id using GRPC Protocol Buffers
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Get a new client on the connection (us)
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	// Finally, convert from proto format to JSON and return
	return model.MetadataFromProto(resp.Metadata), nil
}
