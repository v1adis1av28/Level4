package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"mygrep/internal/config"
	"mygrep/internal/network"
)

func main() {
	conf, err := config.ParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if conf.NodeType != config.NodeTypeCoordinator {
		fmt.Fprintln(os.Stderr, "This is coordinator mode only")
		os.Exit(1)
	}

	log.Printf("Coordinator started with workers: %v", conf.Workers)

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Read %d lines", len(lines))

	parts := splitLines(lines, len(conf.Workers))
	log.Printf("Split into %d parts", len(parts))

	reqFlags := network.GrepFlags{
		Ignore:      conf.Flags.IgnoreFlag,
		Invert:      conf.Flags.InvertFlag,
		Strict:      conf.Flags.StrictStringFlag,
		LineNumbers: conf.Flags.LineNumberFlag,
		CountOnly:   conf.Flags.CountOfLineFlag,
		After:       conf.Flags.AdditionalLineFlag,
		Before:      conf.Flags.PreviousLineFlag,
		Around:      conf.Flags.AroundLineFlag,
	}

	var wg sync.WaitGroup
	resultChan := make(chan *network.GrepResponse, len(conf.Workers))

	for i, addr := range conf.Workers {
		wg.Add(1)
		go func(i int, addr string) {
			defer wg.Done()

			log.Printf("Sending task to %s", addr)

			req := network.GrepRequest{
				Pattern: conf.Pattern,
				Data:    strings.Join(parts[i], "\n"),
				Flags:   reqFlags,
			}

			resp, err := network.SendGrepRequest(addr, req)
			if err != nil {
				log.Printf("Worker %s returned error: %v", addr, err)
				return
			}

			log.Printf("Received response from %s: %+v", addr, *resp)
			resultChan <- resp
		}(i, addr)
	}

	wg.Wait()
	close(resultChan)

	var responses []*network.GrepResponse
	for resp := range resultChan {
		responses = append(responses, resp)
	}

	log.Printf("Collected %d responses", len(responses))

	if len(responses) < conf.Quorum {
		fmt.Fprintf(os.Stderr, "Not enough responses for quorum. Got %d, need %d\n", len(responses), conf.Quorum)
		os.Exit(1)
	}

	for _, resp := range responses {
		for _, line := range resp.Lines {
			fmt.Println(line)
		}
	}
}

func splitLines(lines []string, n int) [][]string {
	if n <= 0 {
		return [][]string{lines}
	}

	chunkSize := len(lines) / n
	if chunkSize == 0 {
		chunkSize = 1
	}

	var result [][]string
	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		result = append(result, lines[i:end])
	}

	if len(result) > n {
		last := result[n-1]
		for i := n; i < len(result); i++ {
			last = append(last, result[i]...)
		}
		result = result[:n]
		result[n-1] = last
	}

	return result
}
