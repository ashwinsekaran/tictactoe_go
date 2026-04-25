package main

import (
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
	l := &lobby{}

	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Printf("connection accepted from %v\n", conn.RemoteAddr())
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
