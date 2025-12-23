package gomoku

import "errors"

var (
	ErrGameOver    = errors.New("игра окончена")
	ErrInvalidTurn = errors.New("недопустимый ход")
)

// Game держит доску и чей ход.
type Game struct {
	Board  *Board
	Turn   Stone
	Winner Stone
	Moves  int
}

// NewGame создаем новую игру с новой доской.
// По умолчанию черные ходят первыми.
func NewGame(size int) *Game {
	return &Game{
		Board:  NewBoard(size),
		Turn:   Black,
		Winner: Empty,
		Moves:  0,
	}
}

// Play ставим фигуру на доску игрока и переключаем ход.
func (g *Game) Play(p Point) error {
	if g.Winner != Empty {
		return ErrGameOver
	}
	if g.Turn != Black && g.Turn != White {
		return ErrInvalidTurn
	}

	if err := g.Board.Set(p, g.Turn); err != nil {
		return err
	}

	g.Moves++

	g.Turn = opposite(g.Turn)
	return nil
}

func opposite(s Stone) Stone {
	if s == Black {
		return White
	}
	return Black
}
