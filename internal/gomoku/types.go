package gomoku

// Фигура на доске
type Stone int

const (
	Empty Stone = iota
	Black
	White
)

func (s Stone) String() string {
	switch s {
	case Black:
		return "●"
	case White:
		return "○"
	default:
		return "."
	}
}

// Point представляет координату на доске.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
