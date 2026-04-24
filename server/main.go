package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("server starting...")

	ser, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("listening on port 8080")

	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Printf("connection accepted from %v\n", conn.RemoteAddr())
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("input: ", line)

		response := "server received: " + line + "\n"
		conn.Write([]byte(response))
	}

	fmt.Printf("client disconnected: %v\n", conn.RemoteAddr())
}
