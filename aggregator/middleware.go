package main

import (
	"github.com/sirupsen/logrus"
	"time"
	"tolling/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"func": "AggregateDistance",
		}).Info("Aggregate distance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)

	return
}

func (l *LogMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": distance,
			"amount":   amount,
		}).Info("Aggregate distance")
	}(time.Now())

	inv, err = l.next.CalculateInvoice(id)

	return
}
