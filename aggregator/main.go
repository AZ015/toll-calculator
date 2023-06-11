package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"tolling/aggregator/client"
	"tolling/types"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "listen adress of HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":50051", "listen adress of GRPC server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)

	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpcListenAddr, svc)
	time.Sleep(time.Second * 5)
	c, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 1.2,
		Unix:  time.Now().UnixNano(),
	})
	if err != nil {
		log.Fatal(err)
	}

	makeHTTPTransport(*httpListenAddr, svc)

	fmt.Println("INVOICER")
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))

	return server.Serve(ln)
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
