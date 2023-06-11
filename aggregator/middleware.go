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
