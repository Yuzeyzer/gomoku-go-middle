package main

import (
	"fmt"

	"github.com/yuzeyzer/gomoku/internal/gomoku"
)

func main() {
	p := gomoku.Point{X: 7, Y: 7}
	fmt.Println("Point:", p)
	fmt.Println("Stone:", gomoku.Black)
}
