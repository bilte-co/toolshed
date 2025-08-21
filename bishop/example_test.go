package bishop_test

import (
	"fmt"
	"github.com/bilte-co/toolshed/bishop"
)

func ExampleGenerate() {
	result := bishop.Generate("hello world")
	fmt.Println(result)
}

func ExampleGenerateFromString_customOptions() {
	opts := &bishop.Options{
		Width:      11,
		Height:     7,
		Symbols:    []rune{' ', '.', 'o', '+', '=', '*', 'B', 'O', 'X'},
		StartChar:  '[',
		EndChar:    ']',
		ShowBorder: true,
	}

	result := bishop.GenerateFromString("custom example", opts)
	fmt.Println(result)
}

func ExampleGenerateFromString_noBorder() {
	opts := &bishop.Options{
		Width:      7,
		Height:     5,
		ShowBorder: false,
	}

	result := bishop.GenerateFromString("no border", opts)
	fmt.Println(result)
}

func ExampleGenerateFromStringSHA256() {
	result := bishop.GenerateFromStringSHA256("sha256 example", nil)
	fmt.Println(result)
}

func ExampleNewBoard() {
	// Create a custom board and walk it manually
	opts := &bishop.Options{
		Width:      9,
		Height:     5,
		Symbols:    []rune{'_', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
		StartChar:  'A',
		EndChar:    'Z',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)

	// Walk with some custom data
	board.Walk([]byte{0x5A, 0xA5}) // Mixed movement pattern

	result := board.Render()
	fmt.Println(result)
}
