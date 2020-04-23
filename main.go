package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//var synchronizerURL = os.Getenv("SYNCHRONIZER_IP")

func main() {
	connection, err := registerWorker()
	if err != nil {
		fmt.Println(err)
		return
	}
	go listen(connection)

	select {}
}

func registerWorker() (*websocket.Conn, error) {

	var dialer websocket.Dialer
	connection, _, err := dialer.Dial("ws://localhost:2216/workers/register/?workertype=cloud", make(http.Header))
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func listen(connection *websocket.Conn) {
	for {
		_, buffer, err := connection.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(string(buffer))
	}
}
