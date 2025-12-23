package gomoku

import "errors"

// AIDifficulty определяет уровень сложности AI.
type AIDifficulty int

const (
	AIEasy AIDifficulty = iota
	AINormal
)

var (
	ErrGameOver    = errors.New("игра окончена")
	ErrInvalidTurn = errors.New("недопустимый ход")
)

// Game хранит состояние игры.
type Game struct {
	Board        *Board
	Turn         Stone
	Winner       Stone
	Moves        int
	WinningLine  []Point // координаты победной линии (5 точек)
	AIDifficulty AIDifficulty
}

// NewGame создаёт новую игру с пустой доской.
// По умолчанию первыми ходят чёрные.
func NewGame(size int) *Game {
	return &Game{
		Board:        NewBoard(size),
		Turn:         Black,
		Winner:       Empty,
		Moves:        0,
		WinningLine:  nil,
		AIDifficulty: AINormal,
	}
}

// Play делает ход текущего игрока в точку p.
// После хода проверяется победа (5 в ряд).
// Если победа найдена — сохраняем WinningLine и игра завершается.
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

	// Проверяем победу только от последнего хода.
	if ok, line := g.findWinningLine(p, g.Turn); ok {
		g.Winner = g.Turn
		g.WinningLine = line
		return nil
	}

	g.Turn = opposite(g.Turn)
	return nil
}

// opposite возвращает противоположный цвет камня.
func opposite(s Stone) Stone {
	if s == Black {
		return White
	}
	return Black
}

// findWinningLine проверяет, приводит ли ход к победе (5 подряд),
// и возвращает координаты победной линии.
func (g *Game) findWinningLine(p Point, s Stone) (bool, []Point) {
	// 4 направления: горизонталь, вертикаль, две диагонали
	dirs := [][2]int{
		{1, 0},
		{0, 1},
		{1, 1},
		{1, -1},
	}

	for _, d := range dirs {
		// Собираем все подряд идущие камни в одной линии (в обе стороны)
		neg := g.collectInDirection(p, s, -d[0], -d[1]) // от p в минус-направление
		pos := g.collectInDirection(p, s, d[0], d[1])   // от p в плюс-направление

		// Полная линия: (neg в обратном порядке) + p + pos
		line := make([]Point, 0, len(neg)+1+len(pos))
		for i := len(neg) - 1; i >= 0; i-- {
			line = append(line, neg[i])
		}
		line = append(line, p)
		line = append(line, pos...)

		if len(line) >= 5 {
			// Для подсветки берём ровно 5 точек вокруг последнего хода.
			start := len(line)/2 - 2
			if start < 0 {
				start = 0
			}
			if start+5 > len(line) {
				start = len(line) - 5
			}
			return true, line[start : start+5]
		}
	}

	return false, nil
}

// collectInDirection собирает подряд идущие точки с камнем цвета s
// в одном направлении от точки p (не включая саму p).
func (g *Game) collectInDirection(p Point, s Stone, dx, dy int) []Point {
	points := make([]Point, 0, 8)
	x, y := p.X+dx, p.Y+dy

	for {
		pp := Point{X: x, Y: y}
		if !g.Board.InBounds(pp) {
			break
		}
		st, err := g.Board.Get(pp)
		if err != nil || st != s {
			break
		}

		points = append(points, pp)
		x += dx
		y += dy
	}

	return points
}

func (g *Game) aiNormal() {
	// 1) Попытка выиграть за White
	if p, ok := g.findBestMove(White); ok {
		_ = g.Play(p)
		return
	}

	// 2) Попытка заблокировать победу Black
	if p, ok := g.findBestMove(Black); ok {
		_ = g.Play(p)
		return
	}

	g.aiFallback()
}

func (g *Game) aiEasy() {
	size := g.Board.Size()

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			p := Point{X: x, Y: y}
			st, _ := g.Board.Get(p)
			if st == Empty {
				_ = g.Play(p)
				return
			}
		}
	}
}

func (g *Game) aiFallback() {
	size := g.Board.Size()
	center := size / 2

	for r := 0; r < size; r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				p := Point{X: center + dx, Y: center + dy}
				if g.Board.InBounds(p) {
					st, _ := g.Board.Get(p)
					if st == Empty {
						_ = g.Play(p)
						return
					}
				}
			}
		}
	}
}

// AIMove типо делает ход за White, если сейчас его очередь.
// Использует простую эвристику: выиграть → заблокировать → сыграть рядом.
// AIMove делает ход за White, если сейчас его очередь.
// Использует простую эвристику: выиграть → заблокировать → fallback.
func (g *Game) AIMove() {
	if g.Winner != Empty || g.Turn != White {
		return
	}

	switch g.AIDifficulty {
	case AIEasy:
		g.aiEasy()
	case AINormal:
		g.aiNormal()
	default:
		g.aiNormal()
	}
}

// findBestMove ищет ход, который даёт победу игроку s.
func (g *Game) findBestMove(s Stone) (Point, bool) {
	size := g.Board.Size()

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			p := Point{X: x, Y: y}
			st, _ := g.Board.Get(p)
			if st != Empty {
				continue
			}

			// временно ставим камень
			g.Board.cells[y][x] = s
			ok, _ := g.findWinningLine(p, s)
			g.Board.cells[y][x] = Empty

			if ok {
				return p, true
			}
		}
	}

	return Point{}, false
}
