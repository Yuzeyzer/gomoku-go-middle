package gomoku

import "errors"

var (
	ErrOutOfBounds  = errors.New("координаты вне границ доски")
	ErrCellOccupied = errors.New("ячейка занята")
)

// Доска Гомоку.
type Board struct {
	size  int       // Размер доски (size x size)
	cells [][]Stone // Ячейки доски
}

// NewBoard создает новую доску.
func NewBoard(size int) *Board {
	cells := make([][]Stone, size)
	for i := range size {
		cells[i] = make([]Stone, size)
	}

	return &Board{
		size:  size,
		cells: cells,
	}
}

// Size возвращает размер доски.
func (b *Board) Size() int {
	return b.size
}

// InBounds проверяет, что координаты находятся внутри доски.
func (b *Board) InBounds(p Point) bool {
	return p.X >= 0 && p.X < b.size && p.Y >= 0 && p.Y < b.size
}

// Get возвращает фигуру на доске.
func (b *Board) Get(p Point) (Stone, error) {
	if !b.InBounds(p) {
		return Empty, ErrOutOfBounds
	}
	return b.cells[p.Y][p.X], nil
}

// Set размещает фигуру на доске.
func (b *Board) Set(p Point, s Stone) error {
	if !b.InBounds(p) {
		return ErrOutOfBounds
	}

	if b.cells[p.Y][p.X] != Empty {
		return ErrCellOccupied
	}

	b.cells[p.Y][p.X] = s
	return nil
}
