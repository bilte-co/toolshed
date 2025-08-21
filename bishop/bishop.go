// Package bishop implements the "drunken bishop" algorithm used by SSH keygen
// to generate ASCII art fingerprints from hash data.
package bishop

import (
	"crypto/md5"
	"crypto/sha256"
	"strings"
)

// DefaultSymbols contains the character mapping used by OpenSSH ssh-keygen
var DefaultSymbols = []rune{' ', '.', 'o', '+', '=', '*', 'B', 'O', 'X', '@', '%', '&', '#', '/', '^'}

// Options configures the bishop walk visualization
type Options struct {
	// Width of the grid (default: 17)
	Width int
	// Height of the grid (default: 9)
	Height int
	// Symbols to use for different visit counts (default: DefaultSymbols)
	Symbols []rune
	// StartChar overrides the start position marker (default: 'S')
	StartChar rune
	// EndChar overrides the end position marker (default: 'E')
	EndChar rune
	// ShowBorder adds a decorative border around the output (default: true)
	ShowBorder bool
}

// DefaultOptions returns the standard ssh-keygen configuration
func DefaultOptions() *Options {
	return &Options{
		Width:      17,
		Height:     9,
		Symbols:    DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: true,
	}
}

// Board represents the bishop's walking grid
type Board struct {
	opts     *Options
	grid     [][]int
	startX   int
	startY   int
	currentX int
	currentY int
}

// NewBoard creates a new board with the given options
func NewBoard(opts *Options) *Board {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Ensure minimum dimensions
	if opts.Width < 3 {
		opts.Width = 3
	}
	if opts.Height < 3 {
		opts.Height = 3
	}

	// Ensure symbols are set
	if len(opts.Symbols) == 0 {
		opts.Symbols = DefaultSymbols
	}

	// Initialize grid
	grid := make([][]int, opts.Height)
	for i := range grid {
		grid[i] = make([]int, opts.Width)
	}

	// Start position is center
	startX := (opts.Width - 1) / 2
	startY := (opts.Height - 1) / 2

	return &Board{
		opts:     opts,
		grid:     grid,
		startX:   startX,
		startY:   startY,
		currentX: startX,
		currentY: startY,
	}
}

// Walk performs the drunken bishop walk using the provided data
func (b *Board) Walk(data []byte) {
	// Process each byte
	for _, byteVal := range data {
		// Extract 4 steps from each byte (2 bits per step)
		for i := 0; i < 4; i++ {
			// Extract 2 bits starting from LSB
			step := (byteVal >> (i * 2)) & 0x03
			b.move(step)
		}
	}
}

// move executes a single step based on 2-bit direction
func (b *Board) move(direction byte) {
	var dx, dy int

	// Decode direction (00=NW, 01=NE, 10=SW, 11=SE)
	switch direction {
	case 0x00: // NW
		dx, dy = -1, -1
	case 0x01: // NE
		dx, dy = 1, -1
	case 0x02: // SW
		dx, dy = -1, 1
	case 0x03: // SE
		dx, dy = 1, 1
	}

	// Calculate new position with clamping
	newX := b.currentX + dx
	newY := b.currentY + dy

	// Clamp to board boundaries
	if newX < 0 {
		newX = 0
	} else if newX >= b.opts.Width {
		newX = b.opts.Width - 1
	}

	if newY < 0 {
		newY = 0
	} else if newY >= b.opts.Height {
		newY = b.opts.Height - 1
	}

	// Update position and increment visit counter
	b.currentX = newX
	b.currentY = newY
	b.grid[newY][newX]++
}

// Render converts the board to ASCII art
func (b *Board) Render() string {
	var result strings.Builder

	if b.opts.ShowBorder {
		// Top border
		result.WriteString("+")
		result.WriteString(strings.Repeat("-", b.opts.Width))
		result.WriteString("+\n")
	}

	// Render each row
	for y := 0; y < b.opts.Height; y++ {
		if b.opts.ShowBorder {
			result.WriteString("|")
		}

		for x := 0; x < b.opts.Width; x++ {
			char := b.getCharForPosition(x, y)
			result.WriteRune(char)
		}

		if b.opts.ShowBorder {
			result.WriteString("|")
		}
		if y < b.opts.Height-1 || b.opts.ShowBorder {
			result.WriteString("\n")
		}
	}

	if b.opts.ShowBorder {
		// Bottom border
		result.WriteString("+")
		result.WriteString(strings.Repeat("-", b.opts.Width))
		result.WriteString("+")
	}

	return result.String()
}

// getCharForPosition returns the appropriate character for a grid position
func (b *Board) getCharForPosition(x, y int) rune {
	// Check for special positions
	if x == b.startX && y == b.startY {
		return b.opts.StartChar
	}
	if x == b.currentX && y == b.currentY && (x != b.startX || y != b.startY) {
		return b.opts.EndChar
	}

	// Get visit count
	visits := b.grid[y][x]

	// Map to symbol
	if len(b.opts.Symbols) == 0 {
		return ' '
	}
	if visits >= len(b.opts.Symbols) {
		visits = len(b.opts.Symbols) - 1
	}

	return b.opts.Symbols[visits]
}

// GenerateFromString creates ASCII art from a string using MD5 hash
func GenerateFromString(input string, opts *Options) string {
	hash := md5.Sum([]byte(input))
	return GenerateFromBytes(hash[:], opts)
}

// GenerateFromStringSHA256 creates ASCII art from a string using SHA256 hash
func GenerateFromStringSHA256(input string, opts *Options) string {
	hash := sha256.Sum256([]byte(input))
	return GenerateFromBytes(hash[:], opts)
}

// GenerateFromBytes creates ASCII art from raw bytes
func GenerateFromBytes(data []byte, opts *Options) string {
	board := NewBoard(opts)
	board.Walk(data)
	return board.Render()
}

// Generate creates ASCII art using default options and MD5
func Generate(input string) string {
	return GenerateFromString(input, nil)
}
