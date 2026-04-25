package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	numGames := flag.Int("games", 1, "number of simultaneous games to simulate")
	flag.Parse()

	server := exec.Command("go", "run", "./server")
	if err := server.Start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	defer server.Process.Kill()

	time.Sleep(500 * time.Millisecond)

	done := make(chan string, *numGames*2)

	fmt.Printf("Starting %d simultaneous games (%d clients)...\n", *numGames, *numGames*2)

	for i := 1; i <= *numGames*2; i++ {
		go simulateClient(fmt.Sprintf("Client-%d", i), done)
	}

	for i := 0; i < *numGames*2; i++ {
		fmt.Println(<-done)
	}

	fmt.Println("All games complete.")
}

func simulateClient(name string, done chan string) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("%s could not connect: %v", name, err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("[%s] %s\n", name, line)

		if strings.Contains(line, "Your turn") {
			time.Sleep(500 * time.Millisecond) // small delay, feels natural
			move := randomMove()
			fmt.Printf("[%s] playing: %s\n", name, move)
			conn.Write([]byte(move + "\n"))
		}

		if strings.Contains(line, "play again") {
			conn.Write([]byte("n\n")) // always say no in tests
		}

		if strings.Contains(line, "You win") || strings.Contains(line, "You lose") || strings.Contains(line, "draw") {
			done <- fmt.Sprintf("%s: game over", name)
			return
		}
	}

	done <- fmt.Sprintf("%s: disconnected", name)
}

func randomMove() string {
	return fmt.Sprintf("%d,%d", rand.Intn(3), rand.Intn(3))
}
