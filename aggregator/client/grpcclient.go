package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tolling/types"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, distance)

	return err
}

func (c *GRPCClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID:         id,
		TotalAmount:   12.45,
		TotalDistance: 123.354,
	}, nil
}
