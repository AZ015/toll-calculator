package main

import (
	"log"
)

type DistanceCalculator struct {
	consumer DataConsumer
}

const kafkaTopic = "obudata"

func main() {
	var (
		svc CalculatorServicer
		err error
	)

	svc = NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
