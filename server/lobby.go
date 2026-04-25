package main

import (
	"net"
	"sync"
)

// lobby holds at most one player waiting for an opponent
type lobby struct {
	mu      sync.Mutex
	waiting net.Conn
}

// join either stores the player as waiting (if lobby is empty) or pairs them
// with the existing waiting player. The mutex prevents a race condition where
// two players joining simultaneously could both see an empty lobby
func (l *lobby) join(conn net.Conn) (opponent net.Conn, paired bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.waiting == nil {
		l.waiting = conn
		return nil, false
	}

	// Second player arrived — grab the waiting player, clear the slot
	opponent = l.waiting
	l.waiting = nil

	return opponent, true
}
