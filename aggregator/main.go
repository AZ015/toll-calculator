package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/invoice", handleGetInvoice(svc))
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

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": "missing OBU ID"})

			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": "invalid OBU ID"})

			return
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": err.Error()})

			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"invoice": invoice})
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(rw).Encode(v)
}
