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
		// Accept blocks until a client connects — main goroutine never handles game logic,
		// it only accepts and hands off connections
		conn, err := ser.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Printf("Connection accepted from %v\n", conn.RemoteAddr())

		// Keepalive detects silent network drops (wifi loss, cable unplug) that TCP
		// wouldn't otherwise notice until a read/write is attempted
		tcpConn := conn.(*net.TCPConn)
		if err := tcpConn.SetKeepAlive(true); err != nil {
			log.Printf("keepalive error for %v: %v", conn.RemoteAddr(), err)
		}
		if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
			log.Printf("keepalive period error for %v: %v", conn.RemoteAddr(), err)
		}

		// Each client gets its own goroutine so Accept() can immediately loop back
		// and handle the next incoming connection without being blocked
		go handleClient(conn, l)
	}
}

// handleClient either puts the player in the lobby to wait, or pairs them with
// a waiting opponent and starts the game.
// connections are closed inside startGame once both players are done
func handleClient(conn net.Conn, l *lobby) {
	if _, err := conn.Write([]byte("Waiting for opponent...\n")); err != nil {
		log.Printf("write error to %v: %v", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	opponent, paired := l.join(conn)

	if !paired {
		// This goroutine returns but conn stays alive — stored inside the lobby.
		// The next player's goroutine will pick it up and start the game
		fmt.Printf("%v is waiting for opponent\n", conn.RemoteAddr())
		return
	}

	fmt.Printf("Pairing %v with %v\n", conn.RemoteAddr(), opponent.RemoteAddr())
	startGame(opponent, conn, l)
}
