package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func run(xConn, oConn net.Conn, l *lobby) {
	var b board

	players := [2]net.Conn{xConn, oConn}
	symbols := [2]string{"X", "O"}

	write := func(msg string) {
		xConn.Write([]byte(msg))
		oConn.Write([]byte(msg))
	}

	write("Game started\n")
	xConn.Write([]byte("You are X\n"))
	oConn.Write([]byte("You are O\n"))
	write(b.display())

	for turn := 0; ; turn = 1 - turn {
		current := players[turn]
		other := players[1-turn]
		symbol := symbols[turn]

		current.Write([]byte(fmt.Sprintf("Your turn (%s) - enter row,col (e.g. 0,0 for first cell):\n", symbol)))
		other.Write([]byte("Waiting for opponent's move...\n"))

		scanner := bufio.NewScanner(current)
		if !scanner.Scan() {
			other.Write([]byte("Opponent disconnected. You win!\n"))
			fmt.Printf("Player %s disconnected\n", symbol)
			go handlePlayAgain(other, l)
			return
		}

		input := strings.TrimSpace(scanner.Text())
		var row, col int

		_, err := fmt.Sscanf(input, "%d,%d", &row, &col)
		if err != nil || row < 0 || row > 2 || col < 0 || col > 2 {
			current.Write([]byte("Invalid input. Please try again (e.g. 1,1)\n"))
			turn = 1 - turn
			continue
		}

		if !b.place(row, col, symbol) {
			current.Write([]byte("Cell is already taken. Please try again.\n"))
			turn = 1 - turn
			continue
		}

		write(b.display())

		if w := b.winner(); w != "" {
			current.Write([]byte("You win!\n"))
			other.Write([]byte("You lose!\n"))
			fmt.Printf("Player %s wins!\n", w)
			break
		}

		if b.isFull() {
			write("It's a draw!\n")
			fmt.Println("Game ended in a draw")
			break
		}

	}

	go handlePlayAgain(xConn, l)
	go handlePlayAgain(oConn, l)
}

func handlePlayAgain(conn net.Conn, l *lobby) {
	conn.Write([]byte("Do you want to play again? (y/n):\n"))

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		conn.Close()
		return
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" {
		conn.Write([]byte("Thanks for playing. Goodbye!\n"))
		conn.Close()
		return
	}

	conn.Write([]byte("Waiting for opponent...\n"))
	opponent, paired := l.join(conn)
	if !paired {
		return
	}

	run(opponent, conn, l)
}
