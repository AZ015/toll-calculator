package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"time"
	"tolling/types"
)

const wsEndpoint = "ws://127.0.0.1:3000/ws"
const sendInterval = time.Second

func sendOBUData(conn *websocket.Conn, data types.OBUData) error {
	return conn.WriteJSON(data)
}

func genLatLong() (float64, float64) {
	return genCord(), genCord()
}

func genCord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", data)

		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
