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
	defer conn.Close()

	conn.Write([]byte("Waiting for opponent...\n"))

	opponent, paired := l.join(conn)

	if !paired {
		fmt.Printf("%v is waiting for opponent\n", conn.RemoteAddr())
		return
	}

	fmt.Printf("Pairing %v with %v\n", conn.RemoteAddr(), opponent.RemoteAddr())
	conn.Write([]byte("Opponent found! Game starting...\n"))
	opponent.Write([]byte("Opponent found! Game starting...\n"))
	//defer conn.Close()
	//
	//scanner := bufio.NewScanner(conn)
	//for scanner.Scan() {
	//	line := scanner.Text()
	//	fmt.Println("input: ", line)
	//
	//	response := "server received: " + line + "\n"
	//	conn.Write([]byte(response))
	//}
	//
	//fmt.Printf("client disconnected: %v\n", conn.RemoteAddr())
}
