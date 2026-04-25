package main

// board is a 3x3 grid — empty cell is "", played cell is "X" or "O"
type board [3][3]string

// display renders the board as a human-readable string for terminal output
func (b *board) display() string {
	result := "\n"
	for i, row := range b {
		for j, column := range row {
			if column == "" {
				result += " _"
			} else {
				result += " " + column
			}
			if j < 2 {
				result += " |"
			}
		}
		result += "\n"
		if i < 2 {
			result += "---|---|---\n"
		}
	}
	return result
}

// place puts a symbol at row,col — returns false if the cell is already taken
func (b *board) place(row, col int, symbol string) bool {
	if b[row][col] != "" {
		return false
	}
	b[row][col] = symbol
	return true
}

// isFull returns true when no empty cells remain — used to detect a draw
func (b *board) isFull() bool {
	for _, row := range b {
		for _, col := range row {
			if col == "" {
				return false
			}
		}
	}
	return true
}

// winner checks all 8 winning lines (3 rows, 3 cols, 2 diagonals)
// returns the winning symbol or "" if no winner yet
func (b *board) winner() string {
	lines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, // top row
		{{1, 0}, {1, 1}, {1, 2}}, // middle row
		{{2, 0}, {2, 1}, {2, 2}}, // bottom row
		{{0, 0}, {1, 0}, {2, 0}}, // left col
		{{0, 1}, {1, 1}, {2, 1}}, // middle col
		{{0, 2}, {1, 2}, {2, 2}}, // right col
		{{0, 0}, {1, 1}, {2, 2}}, // diagonal top-left
		{{0, 2}, {1, 1}, {2, 0}}, // diagonal top-right
	}
	for _, line := range lines {
		a := b[line[0][0]][line[0][1]]
		b2 := b[line[1][0]][line[1][1]]
		c := b[line[2][0]][line[2][1]]

		if a != "" && a == b2 && b2 == c {
			return a
		}
	}
	return ""
}
