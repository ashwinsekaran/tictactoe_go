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

	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fmt.Println("Disconnected from server.")
		close(done)
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := conn.Write([]byte(line + "\n")); err != nil {
			log.Println("write error:", err)
			return
		}

		select {
		case <-done:
			return
		default:
		}
	}
	<-done
}
