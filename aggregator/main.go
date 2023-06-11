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

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)

	svc = NewLogMiddleware(svc)

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
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})

			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})

			return
		}
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(rw).Encode(v)
}
