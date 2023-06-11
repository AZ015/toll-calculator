package main

import (
	"log"
	"tolling/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

func main() {
	var (
		svc CalculatorServicer
		err error
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	//grpcClient, _ := client.NewGRPCClient(aggregatorEndpoint)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
