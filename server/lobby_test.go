package main

import (
	"net"
	"testing"
)

// newPipeConn returns a dummy net.Conn using an in-memory pipe
func newPipeConn() net.Conn {
	client, _ := net.Pipe()
	return client
}

func TestLobbyFirstPlayerWaits(t *testing.T) {
	l := &lobby{}
	conn := newPipeConn()
	defer conn.Close()

	opponent, paired := l.join(conn)

	if paired {
		t.Error("expected first player to wait, got paired=true")
	}
	if opponent != nil {
		t.Error("expected no opponent for first player, got non-nil")
	}
	if l.waiting != conn {
		t.Error("expected first player to be stored in lobby")
	}
}

func TestLobbySecondPlayerGetsPaired(t *testing.T) {
	l := &lobby{}
	conn1 := newPipeConn()
	conn2 := newPipeConn()
	defer conn1.Close()
	defer conn2.Close()

	l.join(conn1)
	opponent, paired := l.join(conn2)

	if !paired {
		t.Error("expected second player to be paired, got paired=false")
	}
	if opponent != conn1 {
		t.Error("expected opponent to be first player")
	}
}

func TestLobbyEmptyAfterPairing(t *testing.T) {
	l := &lobby{}
	conn1 := newPipeConn()
	conn2 := newPipeConn()
	defer conn1.Close()
	defer conn2.Close()

	l.join(conn1)
	l.join(conn2)

	if l.waiting != nil {
		t.Error("expected lobby to be empty after pairing, got non-nil")
	}
}

func TestLobbyThirdPlayerWaitsAgain(t *testing.T) {
	l := &lobby{}
	conn1 := newPipeConn()
	conn2 := newPipeConn()
	conn3 := newPipeConn()
	defer conn1.Close()
	defer conn2.Close()
	defer conn3.Close()

	l.join(conn1)
	l.join(conn2) // paired — lobby cleared
	_, paired := l.join(conn3)

	if paired {
		t.Error("expected third player to wait, got paired=true")
	}
	if l.waiting != conn3 {
		t.Error("expected third player to be stored in lobby")
	}
}

func TestLobbyStartsEmpty(t *testing.T) {
	l := &lobby{}
	if l.waiting != nil {
		t.Error("expected new lobby to be empty")
	}
}
