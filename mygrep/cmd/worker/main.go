package main

import (
	"fmt"
	"log"
	"os"

	"mygrep/internal/config"
	"mygrep/internal/network"
)

func main() {
	conf, err := config.ParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if conf.NodeType == config.NodeTypeWorker {
		log.Printf("Starting worker on port %s", conf.Port)
		err = network.StartWorkerServer(conf.Port)
		if err != nil {
			log.Fatalf("Worker server failed: %v", err)
		}
	} else {
		log.Println("Not running in worker mode")
		os.Exit(1)
	}
}
