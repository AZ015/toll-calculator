package main

import "tolling/types"

type GRPCAggregatorServer struct {
	types.UnimplementedDistanceAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

// transport layer
// JSON -> types.Distance -> all done
// GRPC -> types.AggregateRequest -> type.Distance
// Webpack -> types.Webpack -> type.Distance

// business layer -> business layer type(main type everyone needs to convert)

func (s *GRPCAggregatorServer) AggregateDistance(req *types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return s.svc.AggregateDistance(distance)
}
