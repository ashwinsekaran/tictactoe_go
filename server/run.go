package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func run(xConn, oConn net.Conn, l *lobby) {
	var b board

	players := [2]net.Conn{xConn, oConn}
	symbols := [2]string{"X", "O"}

	writeBoth := func(msg string) {
		if _, err := xConn.Write([]byte(msg)); err != nil {
			log.Printf("write error to X %v: %v", xConn.RemoteAddr(), err)
		}
		if _, err := oConn.Write([]byte(msg)); err != nil {
			log.Printf("write error to O %v: %v", oConn.RemoteAddr(), err)
		}
	}

	writeBoth("Game started\n")

	if _, err := xConn.Write([]byte("You are X\n")); err != nil {
		log.Printf("write error to X %v: %v", xConn.RemoteAddr(), err)
		return
	}
	if _, err := oConn.Write([]byte("You are O\n")); err != nil {
		log.Printf("write error to O %v: %v", oConn.RemoteAddr(), err)
		return
	}

	writeBoth(b.display())

	for turn := 0; ; turn = 1 - turn {
		current := players[turn]
		other := players[1-turn]
		symbol := symbols[turn]

		if _, err := current.Write([]byte(fmt.Sprintf("Your turn (%s) - enter row,col (e.g. 0,0 for first cell):\n", symbol))); err != nil {
			log.Printf("write error to %v: %v", current.RemoteAddr(), err)
			if _, err := other.Write([]byte("Opponent disconnected. You win!\n")); err != nil {
				log.Printf("write error to %v: %v", other.RemoteAddr(), err)
			}
			go handlePlayAgain(other, l)
			return
		}

		if _, err := other.Write([]byte("Waiting for opponent's move...\n")); err != nil {
			log.Printf("write error to %v: %v", other.RemoteAddr(), err)
		}

		scanner := bufio.NewScanner(current)
		if !scanner.Scan() {
			if _, err := other.Write([]byte("Opponent disconnected. You win!\n")); err != nil {
				log.Printf("write error to %v: %v", other.RemoteAddr(), err)
			}
			fmt.Printf("Player %s disconnected\n", symbol)
			go handlePlayAgain(other, l)
			return
		}

		input := strings.TrimSpace(scanner.Text())
		var row, col int

		_, err := fmt.Sscanf(input, "%d,%d", &row, &col)
		if err != nil || row < 0 || row > 2 || col < 0 || col > 2 {
			if _, err := current.Write([]byte("Invalid input. Please try again (e.g. 1,1)\n")); err != nil {
				log.Printf("write error to %v: %v", current.RemoteAddr(), err)
			}
			turn = 1 - turn
			continue
		}

		if !b.place(row, col, symbol) {
			if _, err := current.Write([]byte("Cell is already taken. Please try again.\n")); err != nil {
				log.Printf("write error to %v: %v", current.RemoteAddr(), err)
			}
			turn = 1 - turn
			continue
		}

		writeBoth(b.display())

		if w := b.winner(); w != "" {
			if _, err := current.Write([]byte("You win!\n")); err != nil {
				log.Printf("write error to %v: %v", current.RemoteAddr(), err)
			}
			if _, err := other.Write([]byte("You lose!\n")); err != nil {
				log.Printf("write error to %v: %v", other.RemoteAddr(), err)
			}
			fmt.Printf("Player %s wins!\n", w)
			break
		}

		if b.isFull() {
			writeBoth("It's a draw!\n")
			fmt.Println("Game ended in a draw")
			break
		}
	}

	go handlePlayAgain(xConn, l)
	go handlePlayAgain(oConn, l)
}

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

	opponent, paired := l.join(conn)
	if !paired {
		return
	}

	run(opponent, conn, l)
}
