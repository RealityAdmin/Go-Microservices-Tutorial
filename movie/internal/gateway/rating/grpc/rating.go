package grpc

import (
	"context"

	"movieexamplekhubaib.com/gen"
	"movieexamplekhubaib.com/internal/grpcutil"
	"movieexamplekhubaib.com/pkg/discovery"
	"movieexamplekhubaib.com/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get aggregated ratings for a record using GRPC
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(recordType),
	})

	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil
}
