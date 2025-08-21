package bishop_test

import (
	"strings"
	"testing"

	"github.com/bilte-co/toolshed/bishop"
	"github.com/stretchr/testify/require"
)

func TestDefaultOptions(t *testing.T) {
	opts := bishop.DefaultOptions()
	require.Equal(t, 17, opts.Width)
	require.Equal(t, 9, opts.Height)
	require.Equal(t, bishop.DefaultSymbols, opts.Symbols)
	require.Equal(t, 'S', opts.StartChar)
	require.Equal(t, 'E', opts.EndChar)
	require.True(t, opts.ShowBorder)
}

func TestNewBoard_DefaultOptions(t *testing.T) {
	board := bishop.NewBoard(nil)
	result := board.Render()

	// Should have border and center start position
	lines := strings.Split(result, "\n")
	require.Greater(t, len(lines), 3) // At least top border, content, bottom border
	require.True(t, strings.HasPrefix(lines[0], "+"))
	require.True(t, strings.HasSuffix(lines[0], "+"))
}

func TestNewBoard_CustomDimensions(t *testing.T) {
	opts := &bishop.Options{
		Width:      7,
		Height:     5,
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: true,
	}

	board := bishop.NewBoard(opts)
	result := board.Render()

	lines := strings.Split(result, "\n")
	// Should have 5 content lines + 2 border lines = 7 total
	require.Equal(t, 7, len(lines))

	// Check width (7 content + 2 border chars = 9)
	require.Equal(t, 9, len(lines[0])) // "+-------+"
	require.Equal(t, 9, len(lines[1])) // "|content|"
}

func TestNewBoard_MinimumDimensions(t *testing.T) {
	opts := &bishop.Options{
		Width:      1, // Should be clamped to 3
		Height:     2, // Should be clamped to 3
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)
	result := board.Render()

	lines := strings.Split(result, "\n")
	require.Equal(t, 3, len(lines))    // 3 rows minimum
	require.Equal(t, 3, len(lines[0])) // 3 columns minimum
}

func TestBoardWalk_SimpleInput(t *testing.T) {
	opts := &bishop.Options{
		Width:      5,
		Height:     3,
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)

	// Use simple input: single byte with value 0x00 (all NW moves)
	board.Walk([]byte{0x00})

	result := board.Render()

	// Start should be at center (2,1), after 4 NW moves should be at (0,0)
	// But clamped, so likely at top-left corner
	require.Contains(t, result, "E") // End marker should be present
	require.Contains(t, result, "S") // Start marker should be present
}

func TestBoardWalk_AllDirections(t *testing.T) {
	opts := &bishop.Options{
		Width:      7,
		Height:     7,
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)

	// Test all 4 directions: 0x00 (NW), 0x55 (NE), 0xAA (SW), 0xFF (SE)
	board.Walk([]byte{0x00, 0x55, 0xAA, 0xFF})

	result := board.Render()

	// Should contain start and end markers
	require.Contains(t, result, "S")
	require.Contains(t, result, "E")

	// Should contain some visit symbols (not just spaces)
	require.True(t, strings.Contains(result, ".") || strings.Contains(result, "o"))
}

func TestCustomSymbols(t *testing.T) {
	customSymbols := []rune{'_', '1', '2', '3', '4', '5'}
	opts := &bishop.Options{
		Width:      5,
		Height:     3,
		Symbols:    customSymbols,
		StartChar:  'A',
		EndChar:    'Z',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)
	board.Walk([]byte{0x00}) // Some movement

	result := board.Render()

	// Should use custom symbols
	require.Contains(t, result, "A")    // Custom start
	require.Contains(t, result, "Z")    // Custom end
	require.NotContains(t, result, "S") // Default start should not appear
	require.NotContains(t, result, "E") // Default end should not appear
}

func TestGenerateFromString_Consistency(t *testing.T) {
	input := "test"

	// Same input should produce same output
	result1 := bishop.GenerateFromString(input, nil)
	result2 := bishop.GenerateFromString(input, nil)

	require.Equal(t, result1, result2)

	// Different inputs should produce different outputs
	result3 := bishop.GenerateFromString("different", nil)
	require.NotEqual(t, result1, result3)
}

func TestGenerateFromStringSHA256(t *testing.T) {
	input := "test"

	resultMD5 := bishop.GenerateFromString(input, nil)
	resultSHA256 := bishop.GenerateFromStringSHA256(input, nil)

	// Should be different since they use different hashes
	require.NotEqual(t, resultMD5, resultSHA256)

	// But both should contain borders and markers
	require.Contains(t, resultSHA256, "+")
	require.Contains(t, resultSHA256, "S")
	require.Contains(t, resultSHA256, "E")
}

func TestGenerate_ConvenienceFunction(t *testing.T) {
	result := bishop.Generate("test")

	// Should produce valid output with defaults
	require.Contains(t, result, "+")
	require.Contains(t, result, "S")
	require.Contains(t, result, "E")

	lines := strings.Split(result, "\n")
	require.Greater(t, len(lines), 5) // Should have multiple lines
}

func TestNoBorder(t *testing.T) {
	opts := &bishop.Options{
		Width:      5,
		Height:     3,
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: false,
	}

	result := bishop.GenerateFromString("test", opts)

	// Should not contain border characters
	require.NotContains(t, result, "+")
	require.NotContains(t, result, "|")

	// But should contain content
	require.Contains(t, result, "S")
	require.Contains(t, result, "E")

	lines := strings.Split(result, "\n")
	require.Equal(t, 3, len(lines))    // Exact height, no border
	require.Equal(t, 5, len(lines[0])) // Exact width, no border
}

func TestSameStartEnd(t *testing.T) {
	// Test case where bishop doesn't move (stays at start)
	opts := &bishop.Options{
		Width:      3,
		Height:     3,
		Symbols:    bishop.DefaultSymbols,
		StartChar:  'S',
		EndChar:    'E',
		ShowBorder: false,
	}

	board := bishop.NewBoard(opts)
	// Don't walk anywhere - stay at start position
	result := board.Render()

	// When start and end are same, should show start marker
	require.Contains(t, result, "S")
	require.NotContains(t, result, "E") // End shouldn't override start when they're the same
}

func TestLargeInput(t *testing.T) {
	// Test with larger input to ensure algorithm handles more steps
	largeInput := strings.Repeat("test data for longer walks", 10)

	result := bishop.Generate(largeInput)

	require.Contains(t, result, "S")
	require.Contains(t, result, "E")
	require.Contains(t, result, "+")

	// Should have some visited cells (not just start/end)
	hasVisitedCells := false
	for _, char := range bishop.DefaultSymbols[1:] { // Skip space (index 0)
		if strings.ContainsRune(result, char) {
			hasVisitedCells = true
			break
		}
	}
	require.True(t, hasVisitedCells, "Should have some visited cells with symbols")
}
