package main

import (
	"net"
	"sync"
)

type lobby struct {
	mu      sync.Mutex
	waiting net.Conn
}

func (l *lobby) join(conn net.Conn) (opponent net.Conn, paired bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.waiting == nil {
		l.waiting = conn
		return nil, false
	}
	opponent = l.waiting
	l.waiting = nil

	return opponent, true
}
