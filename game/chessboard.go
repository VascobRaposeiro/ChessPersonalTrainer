package game

import (
	"fmt"
	"math/rand/v2"
)

// Casa do tabuleiro
type Square struct {
	Column int // 1-8 (A-H)
	Row    int // 1-8
}

// Gerar Coordenada
func GenerateRandomSquare() Square {

	var col int
	var row int

	col = rand.IntN(8) + 1
	row = rand.IntN(8) + 1

	sq := Square{}
	sq.Column = col
	sq.Row = row

	return sq
}

// Verificar se a casa e preta ou branca
func (s Square) IsBlack() bool {

	var blackSquare bool

	if (s.Column+s.Row)%2 == 0 {
		blackSquare = true
	} else {
		blackSquare = false
	}

	return blackSquare
}

// Converter coordenada para string
func (s Square) String() string {

	coordinate := [8]string{"A", "B", "C", "D", "E", "F", "G", "H"}
	var col = coordinate[s.Column-1]
	var row = s.Row
	str := fmt.Sprint(col, row)

	return str
}
