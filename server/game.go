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
