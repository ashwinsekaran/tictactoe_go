package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("Server starting...")

	ser, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Listening on port 8080")
	l := &lobby{}

	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Printf("Connection accepted from %v\n", conn.RemoteAddr())

		tcpConn := conn.(*net.TCPConn)
		if err := tcpConn.SetKeepAlive(true); err != nil {
			log.Printf("keepalive error for %v: %v", conn.RemoteAddr(), err)
		}
		if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
			log.Printf("keepalive period error for %v: %v", conn.RemoteAddr(), err)
		}

		go handleClient(conn, l)
	}

}

func handleClient(conn net.Conn, l *lobby) {
	if _, err := conn.Write([]byte("Waiting for opponent...\n")); err != nil {
		log.Printf("write error to %v: %v", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	opponent, paired := l.join(conn)

	if !paired {
		fmt.Printf("%v is waiting for opponent\n", conn.RemoteAddr())
		return
	}

	fmt.Printf("Pairing %v with %v\n", conn.RemoteAddr(), opponent.RemoteAddr())
	run(opponent, conn, l)
}
