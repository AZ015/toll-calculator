package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"tolling/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen adress of HTTP server")
	flag.Parse()

	store := NewMemoryStore()

	var (
		svc = NewInvoiceAggregator(store)
	)

	makeHTTPTransport(*listenAddr, svc)

	fmt.Println("INVOICER")
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}
}
