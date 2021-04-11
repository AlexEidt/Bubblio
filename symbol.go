package main

type symbol struct {
	width  int
	height int
	bitmap [][]bool
}

// Turns a bitmap into a struct used to draw letters.
func NewSymbol(data []string, identifier rune) *symbol {
	bitmap := make([][]bool, len(data))
	for row := 0; row < len(data); row++ {
		bitmap[row] = make([]bool, len(data[row]))
		for column, letter := range data[row] {
			bitmap[row][column] = letter == identifier
		}
	}
	return &symbol{
		width:  len(data[0]),
		height: len(data),
		bitmap: bitmap,
	}
}
