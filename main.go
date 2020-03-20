package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//var synchronizerURL = os.Getenv("SYNCHRONIZER_IP")

func main() {
	uuid := registerWorker()
	fmt.Println(uuid)
}

func registerWorker() string {

	// Get ip address from env
	ip := os.Getenv("CLOUDWORKER_IP")

	// Build http post request to /workers endpoint
	reqBody, _ := json.Marshal(map[string]interface{}{"ip": ip, "workerType": "cloud_worker"})
	req, err := http.NewRequest("POST", "http://localhost:2216/workers/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic("Failed to create request: " + err.Error())
	}
	fmt.Printf("%v\n", req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic("Failed to register worker: " + err.Error())
	}
	uuid, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.Status != "200 OK" {
		panic("Failed to register worker: " + res.Status + " " + string(uuid))
	}

	return string(uuid)
}
