package main

type board [3][3]string

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
func (b *board) place(row, col int, symbol string) bool {
	if b[row][col] != "" {
		return false
	}
	b[row][col] = symbol
	return true
}

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

func (b *board) winner() string {
	lines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
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
