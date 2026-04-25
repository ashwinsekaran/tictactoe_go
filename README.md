# Tic-Tac-Toe over TCP (Go)

A networked Tic-Tac-Toe game written in Go. Two players connect to a server over TCP and play against each other in the terminal.

---

## Project Structure

```
tictactoe_go/
├── server/
│   ├── main.go        — entry point, accepts TCP connections, spawns a goroutine per client
│   ├── lobby.go       — pairs two waiting players together using a mutex-protected waiting room
│   ├── game.go        — board representation, move placement, win/draw detection
│   ├── run.go         — starts game, manages turns, handles play again flow
│   ├── game_test.go   — unit tests for board logic (place, winner, isFull)
│   └── lobby_test.go  — unit tests for lobby pairing logic
├── client/
│   └── main.go     — connects to server, two goroutines for reading server messages and stdin
└── simulate/
    └── main.go     — automated simulation script, simulates N simultaneous games with bot clients
```

---

## Requirements

- Go 1.21 or higher

---

## How to Run

### 1. Start the Server

```bash
go run ./server
```

Server listens on port 8080. Keep this running in its own terminal.

### 2. Start Clients (one per player)

Open a new terminal for each player:

```bash
go run ./client
```

- First client waits for an opponent
- Second client connects and the game begins
- Players take turns entering moves as `row,col` (e.g. `0,0` for top-left, `2,2` for bottom-right)

```
Board positions:
 0,0 | 0,1 | 0,2
-----|-----|-----
 1,0 | 1,1 | 1,2
-----|-----|-----
 2,0 | 2,1 | 2,2
```

### 3. Play Again

After a game ends, each player is asked `Do you want to play again? (y/n)`:
- `y` — puts the player back in the lobby to wait for any available opponent
- `n` — disconnects the player

---

## Running Multiple Games Simultaneously

You can run 1, 10, or 100 games at the same time by opening more client terminals. The server pairs clients in the order they connect — every two clients form one game.

---

## Unit Tests

```bash
# run all unit tests
go test ./server

# run with details
go test ./server -v

# run a specific test
go test ./server -run TestWinnerTopRow
```

---

## Automated Simulation

The test script starts the server automatically and simulates bot clients playing full games.

```bash
# run 1 game (2 bot clients)
go run ./simulate -games 1

# run 10 simultaneous games (20 bot clients)
go run ./simulate -games 10

# run 100 simultaneous games (200 bot clients)
go run ./simulate -games 100

# see all options
go run ./simulate -help
```

No need to start the server manually — the test script handles it.

---

## Design

- **Transport**: raw TCP using Go's `net` package — no external dependencies
- **Message framing**: newline-delimited plain text
- **Concurrency**: one goroutine per connected client on the server; game state is owned by a single goroutine per game — no shared mutable state between games
- **Pairing**: mutex-protected lobby queues players and pairs them as they arrive
- **Client**: two goroutines — one reads from the server, one reads from stdin — allowing the player to see opponent moves immediately while waiting to type
