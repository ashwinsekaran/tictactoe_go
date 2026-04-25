package main

import "testing"

// ---- place() tests ----

func TestPlaceOnEmptyCell(t *testing.T) {
	var b board
	ok := b.place(0, 0, "X")
	if !ok {
		t.Error("expected true for empty cell, got false")
	}
	if b[0][0] != "X" {
		t.Errorf("expected X at 0,0, got %s", b[0][0])
	}
}

func TestPlaceOnOccupiedCell(t *testing.T) {
	var b board
	b.place(1, 1, "X")
	ok := b.place(1, 1, "O")
	if ok {
		t.Error("expected false for occupied cell, got true")
	}
	if b[1][1] != "X" {
		t.Errorf("expected X to remain at 1,1, got %s", b[1][1])
	}
}

func TestPlaceAllCells(t *testing.T) {
	var b board
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			ok := b.place(row, col, "X")
			if !ok {
				t.Errorf("expected true for empty cell at %d,%d, got false", row, col)
			}
		}
	}
}

// ---- winner() tests ----

func TestWinnerTopRow(t *testing.T) {
	var b board
	b[0][0], b[0][1], b[0][2] = "X", "X", "X"
	if w := b.winner(); w != "X" {
		t.Errorf("expected X to win top row, got %s", w)
	}
}

func TestWinnerMiddleRow(t *testing.T) {
	var b board
	b[1][0], b[1][1], b[1][2] = "O", "O", "O"
	if w := b.winner(); w != "O" {
		t.Errorf("expected O to win middle row, got %s", w)
	}
}

func TestWinnerBottomRow(t *testing.T) {
	var b board
	b[2][0], b[2][1], b[2][2] = "X", "X", "X"
	if w := b.winner(); w != "X" {
		t.Errorf("expected X to win bottom row, got %s", w)
	}
}

func TestWinnerLeftColumn(t *testing.T) {
	var b board
	b[0][0], b[1][0], b[2][0] = "X", "X", "X"
	if w := b.winner(); w != "X" {
		t.Errorf("expected X to win left column, got %s", w)
	}
}

func TestWinnerMiddleColumn(t *testing.T) {
	var b board
	b[0][1], b[1][1], b[2][1] = "O", "O", "O"
	if w := b.winner(); w != "O" {
		t.Errorf("expected O to win middle column, got %s", w)
	}
}

func TestWinnerRightColumn(t *testing.T) {
	var b board
	b[0][2], b[1][2], b[2][2] = "X", "X", "X"
	if w := b.winner(); w != "X" {
		t.Errorf("expected X to win right column, got %s", w)
	}
}

func TestWinnerDiagonalTopLeft(t *testing.T) {
	var b board
	b[0][0], b[1][1], b[2][2] = "X", "X", "X"
	if w := b.winner(); w != "X" {
		t.Errorf("expected X to win diagonal top-left, got %s", w)
	}
}

func TestWinnerDiagonalTopRight(t *testing.T) {
	var b board
	b[0][2], b[1][1], b[2][0] = "O", "O", "O"
	if w := b.winner(); w != "O" {
		t.Errorf("expected O to win diagonal top-right, got %s", w)
	}
}

func TestNoWinnerEmptyBoard(t *testing.T) {
	var b board
	if w := b.winner(); w != "" {
		t.Errorf("expected no winner on empty board, got %s", w)
	}
}

func TestNoWinnerMidGame(t *testing.T) {
	var b board
	b[0][0], b[1][1] = "X", "O"
	if w := b.winner(); w != "" {
		t.Errorf("expected no winner mid-game, got %s", w)
	}
}

func TestNoWinnerMixedBoard(t *testing.T) {
	var b board
	// draw board — no winner
	b[0][0], b[0][1], b[0][2] = "X", "O", "X"
	b[1][0], b[1][1], b[1][2] = "O", "X", "O"
	b[2][0], b[2][1], b[2][2] = "O", "X", "O"
	if w := b.winner(); w != "" {
		t.Errorf("expected no winner on draw board, got %s", w)
	}
}

// ---- isFull() tests ----

func TestIsFullEmptyBoard(t *testing.T) {
	var b board
	if b.isFull() {
		t.Error("expected false for empty board, got true")
	}
}

func TestIsFullPartialBoard(t *testing.T) {
	var b board
	b[0][0], b[1][1] = "X", "O"
	if b.isFull() {
		t.Error("expected false for partial board, got true")
	}
}

func TestIsFullFullBoard(t *testing.T) {
	var b board
	b[0][0], b[0][1], b[0][2] = "X", "O", "X"
	b[1][0], b[1][1], b[1][2] = "O", "X", "O"
	b[2][0], b[2][1], b[2][2] = "O", "X", "O"
	if !b.isFull() {
		t.Error("expected true for full board, got false")
	}
}

func TestIsFullOneEmpty(t *testing.T) {
	var b board
	b[0][0], b[0][1], b[0][2] = "X", "O", "X"
	b[1][0], b[1][1], b[1][2] = "O", "X", "O"
	b[2][0], b[2][1] = "O", "X"
	// b[2][2] is empty
	if b.isFull() {
		t.Error("expected false when one cell is empty, got true")
	}
}
