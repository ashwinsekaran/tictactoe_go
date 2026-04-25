package main

import (
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
	l := &lobby{}

	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		fmt.Printf("Connection accepted from %v\n", conn.RemoteAddr())
		go handleClient(conn, l)
	}

}

func handleClient(conn net.Conn, l *lobby) {
	conn.Write([]byte("Waiting for opponent...\n"))

	opponent, paired := l.join(conn)

	if !paired {
		fmt.Printf("%v is waiting for opponent\n", conn.RemoteAddr())
		return
	}

	fmt.Printf("Pairing %v with %v\n", conn.RemoteAddr(), opponent.RemoteAddr())
	run(opponent, conn, l)
}
