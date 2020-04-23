package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var connection *websocket.Conn
var tasks chan *Payload = make(chan *Payload, 1024)

type TaskMessage struct {
	MessageType string
	Payload     Payload
}

type Payload struct {
	IntentType string
	TaskType   string
	TaskUUID   string
	Config     interface{}
	Input      Input
}

type Input struct {
	Board []int8 `json:"board"`
	Size  int    `json:"size"`
}

func main() {
	connection, err := registerWorker()
	if err != nil {
		fmt.Println(err)
		return
	}
	go listen(connection)
	go execute(tasks, connection)

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

		// Try to read message into JSON
		//		var message map[string]interface{}
		var message TaskMessage
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(message)

		if message.MessageType == "Intent" {
			payload := message.Payload

			tasks <- &payload
		}
	}
}

func sendMessage(message *map[string]interface{}, connection *websocket.Conn) {
	go func() {
		// Convert to buffer
		buffer, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
		}

		connection.WriteMessage(websocket.TextMessage, buffer)
	}()
}

func execute(tasks chan *Payload, connection *websocket.Conn) {
	for {
		task := <-tasks
		fmt.Println("hello")

		if task.TaskType == "GOL" {
			sendMessage(&map[string]interface{}{
				"MessageType": "WorkResponse",
				"Output":      gol(task.Input.Size, task.Input.Board),
				"start":       1,
				"end":         2,
			}, connection)

		} else {
			fmt.Println("Unknown task")
		}
	}
}

func gol(size int, board []int8) []int8 {
	result := make([]int8, size*size)
	for i := 0; i < size*size; i++ {
		x := i % size
		y := i / size

		x1 := (x + size - 1) % size
		x2 := x
		x3 := (x + size + 1) % size

		y1 := ((y + size - 1) % size) * size
		y2 := y * size
		y3 := ((y + size + 1) % size) * size

		var count int8 = 0

		count += board[x1+y1]
		count += board[x2+y1]
		count += board[x3+y1]

		count += board[x1+y2]
		count += board[x3+y2]

		count += board[x1+y3]
		count += board[x2+y3]
		count += board[x3+y3]

		if board[x2+y2] == 1 {
			if count < 2 {
				result[x2+y2] = 0
			} else if count > 3 {
				result[x2+y2] = 0
			} else {
				result[x2+y2] = 1
			}
		} else {
			if count == 3 {
				result[x2+y2] = 1
			} else {
				result[x2+y2] = 0
			}
		}
	}
	return result
}
