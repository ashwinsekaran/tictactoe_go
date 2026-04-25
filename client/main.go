package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// done is closed when the server disconnects — signals the stdin goroutine to stop
	done := make(chan struct{})

	// Goroutine 1: reads messages from the server and prints them.
	// Runs concurrently with stdin reading so opponent moves appear immediately
	// even while the player is mid-typing
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fmt.Println("Disconnected from server.")
		close(done)
	}()

	// Goroutine 2 (main): reads player input from stdin and sends to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := conn.Write([]byte(line + "\n")); err != nil {
			log.Println("write error:", err)
			return
		}

		// Check if server disconnected between moves
		select {
		case <-done:
			return
		default:
		}
	}

	// Wait for server goroutine to finish before exiting
	<-done
}
