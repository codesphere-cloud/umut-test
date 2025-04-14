package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "http://localhost:3000/"
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		fmt.Printf("failed to construct request: %s", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to SSE stream: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	log.Println("Alles ok.")

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from SSE stream: %v\n", err)
			break
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "event:") {
			eventType := strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			fmt.Printf("Received event type: %s\n", eventType)
		} else if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			fmt.Printf("Received data: %s\n", data)
		} else if line == "" {
			fmt.Println("--- End of Event ---")
		} else if strings.HasPrefix(line, ":") {
			fmt.Printf("Received comment: %s\n", line)
		} else {
			if line != "" {
				fmt.Printf("Received unknown line: %s\n", line)
			}
		}
	}

	fmt.Println("SSE stream closed.")
}
