package main

import (
	"fmt"

	"github.com/yuzeyzer/gomoku/internal/gomoku"
)

func main() {
	// Пример использования доски Гомоку, дожно вернуть "●"
	b := gomoku.NewBoard(15)
	_ = b.Set(gomoku.Point{X: 7, Y: 7}, gomoku.Black)
	stone, _ := b.Get(gomoku.Point{X: 7, Y: 7})
	fmt.Println(stone)
}
