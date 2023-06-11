package client

import (
	"context"
	"tolling/types"
)

type Client interface {
	Aggregate(ctx context.Context, distance *types.AggregateRequest) error
}
