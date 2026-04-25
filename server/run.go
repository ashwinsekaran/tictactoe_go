package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// startGame runs a full game between two players. xConn is always X and goes first.
// This function owns both connections for its lifetime — no other goroutine touches them.
// Connections are not closed here; handlePlayAgain decides whether to close or reuse them.
func startGame(xConn, oConn net.Conn, l *lobby) {
	var b board

	players := [2]net.Conn{xConn, oConn}
	symbols := [2]string{"X", "O"}

	// returns false if either write fails — caller should stop the game
	writeBothPlayers := func(msg string) bool {
		if _, err := xConn.Write([]byte(msg)); err != nil {
			log.Printf("write error to X %v: %v", xConn.RemoteAddr(), err)
			return false
		}
		if _, err := oConn.Write([]byte(msg)); err != nil {
			log.Printf("write error to O %v: %v", oConn.RemoteAddr(), err)
			return false
		}
		return true
	}

	if !writeBothPlayers("Game started\n") {
		return
	}

	if _, err := xConn.Write([]byte("You are X\n")); err != nil {
		log.Printf("write error to X %v: %v", xConn.RemoteAddr(), err)
		return
	}
	if _, err := oConn.Write([]byte("You are O\n")); err != nil {
		log.Printf("write error to O %v: %v", oConn.RemoteAddr(), err)
		return
	}

	if !writeBothPlayers(b.display()) {
		return
	}

	// scanners created once per connection and reused across turns — avoids losing
	// buffered data that would occur if a new scanner was created each iteration
	scanners := [2]*bufio.Scanner{bufio.NewScanner(xConn), bufio.NewScanner(oConn)}

	// turn alternates between 0 and 1 using 1-turn trick: 0→1→0→1...
	for turn := 0; ; turn = 1 - turn {
		player1 := players[turn]
		player2 := players[1-turn]
		symbol := symbols[turn]

		if _, err := player1.Write([]byte(fmt.Sprintf("Your turn (%s) - enter row,col (e.g. 0,0 for first cell):\n", symbol))); err != nil {
			// Write failed means player1 disconnected — notify player2 and exit
			log.Printf("write error to %v: %v", player1.RemoteAddr(), err)
			if _, err := player2.Write([]byte("Opponent disconnected. You win!\n")); err != nil {
				log.Printf("write error to %v: %v", player2.RemoteAddr(), err)
			}
			go handlePlayAgain(player2, l)
			return
		}

		if _, err := player2.Write([]byte("Waiting for opponent's move...\n")); err != nil {
			log.Printf("write error to %v: %v", player2.RemoteAddr(), err)
			go handlePlayAgain(player1, l)
			return
		}

		// Block here until player1 sends a move — scanner.Scan() returns false on disconnect
		if !scanners[turn].Scan() {
			if _, err := player2.Write([]byte("Opponent disconnected. You win!\n")); err != nil {
				log.Printf("write error to %v: %v", player2.RemoteAddr(), err)
			}
			fmt.Printf("Player %s disconnected\n", symbol)
			go handlePlayAgain(player2, l)
			return
		}

		input := strings.TrimSpace(scanners[turn].Text())
		var row, col int

		_, err := fmt.Sscanf(input, "%d,%d", &row, &col)
		if err != nil || row < 0 || row > 2 || col < 0 || col > 2 {
			if _, err := player1.Write([]byte("Invalid input. Please try again (e.g. 1,1)\n")); err != nil {
				log.Printf("write error to %v: %v", player1.RemoteAddr(), err)
			}
			turn = 1 - turn // undo turn flip so same player goes again
			continue
		}

		if !b.place(row, col, symbol) {
			if _, err := player1.Write([]byte("Cell is already taken. Please try again.\n")); err != nil {
				log.Printf("write error to %v: %v", player1.RemoteAddr(), err)
			}
			turn = 1 - turn // undo turn flip so same player goes again
			continue
		}

		if !writeBothPlayers(b.display()) {
			return
		}

		if winner := b.winner(); winner != "" {
			if _, err := player1.Write([]byte("You win!\n")); err != nil {
				log.Printf("write error to %v: %v", player1.RemoteAddr(), err)
			}
			if _, err := player2.Write([]byte("You lose!\n")); err != nil {
				log.Printf("write error to %v: %v", player2.RemoteAddr(), err)
			}
			fmt.Printf("Player %s wins!\n", winner)
			break // break exits the loop, allowing play again logic below to run
		}

		if b.isFull() {
			if !writeBothPlayers("It's a draw!\n") {
				return
			}
			fmt.Println("Game ended in a draw")
			break
		}
	}

	// Each player independently decides whether to play again —
	// no synchronisation needed, lobby handles re-pairing naturally
	go handlePlayAgain(xConn, l)
	go handlePlayAgain(oConn, l)
}

// handlePlayAgain asks one player if they want to play again.
// If yes, puts them back in the lobby — they may be paired with any waiting player.
// If no, closes their connection.
func handlePlayAgain(conn net.Conn, l *lobby) {
	if _, err := conn.Write([]byte("Do you want to play again? (y/n):\n")); err != nil {
		log.Printf("write error to %v: %v", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		conn.Close()
		return
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" {
		if _, err := conn.Write([]byte("Thanks for playing. Goodbye!\n")); err != nil {
			log.Printf("write error to %v: %v", conn.RemoteAddr(), err)
		}
		conn.Close()
		return
	}

	if _, err := conn.Write([]byte("Waiting for opponent...\n")); err != nil {
		log.Printf("write error to %v: %v", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	// Re-enter the lobby — if another player is already waiting, game starts immediately
	opponent, paired := l.join(conn)
	if !paired {
		return
	}

	startGame(opponent, conn, l)
}
