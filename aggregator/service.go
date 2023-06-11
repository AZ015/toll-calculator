package main

import (
	"fmt"
	"tolling/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(distance types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting to storage:", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	fmt.Println("get distance for obu id:", obuID)

	dist, err := i.store.Get(obuID)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return inv, nil
}
