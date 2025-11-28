package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendGrepRequest(workerAddr string, req GrepRequest) (*GrepResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling request for %s: %v", workerAddr, err)
		return nil, err
	}

	resp, err := http.Post("http://"+workerAddr+"/grep", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending request to %s: %v", workerAddr, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from %s: %v", workerAddr, err)
		return nil, err
	}

	var greqResp GrepResponse
	err = json.Unmarshal(body, &greqResp)
	if err != nil {
		log.Printf("Error unmarshalling response from %s: %v", workerAddr, err)
		return nil, err
	}

	if greqResp.Error != "" {
		log.Printf("Worker %s returned error: %s", workerAddr, greqResp.Error)
		return nil, fmt.Errorf("worker error: %s", greqResp.Error)
	}

	log.Printf("Worker %s responded: %+v", workerAddr, greqResp)
	return &greqResp, nil
}
