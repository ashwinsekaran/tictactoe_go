package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Server starting...")

	ser, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Listening on port 8080")
	conn, err := ser.Accept()
	if err != nil {
		log.Fatalf("Failed to accept: %v", err)
	}
	fmt.Printf("connection accepted from %v", conn.RemoteAddr())
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("input: ", line)

		response := "Server received: " + line + "\n"
		conn.Write([]byte(response))
	}
}
