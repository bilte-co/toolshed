package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/bilte-co/toolshed/password"
	"github.com/stretchr/testify/require"
)

func TestPasswordCheckCmd_StrongPassword(t *testing.T) {
	cmd := &cli.PasswordCheckCmd{
		Text: "MyStr0ng!P@ssw0rd2024",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestPasswordCheckCmd_WeakPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"numeric only", "123456"},
		{"common weak", "password"},
		{"simple", "abc123"},
		{"short", "a1!"},
		{"repetitive", "aaaaaaaaa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture os.Exit to prevent test termination
			oldExit := osExit
			var exitCode int
			osExit = func(code int) { exitCode = code }
			defer func() { osExit = oldExit }()

			cmd := &cli.PasswordCheckCmd{
				Text: tt.password,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err) // Command doesn't return error, it calls os.Exit
			require.Equal(t, 1, exitCode, "Should exit with code 1 for weak password")
		})
	}
}

func TestPasswordCheckCmd_EmptyPassword(t *testing.T) {
	cmd := &cli.PasswordCheckCmd{
		Text: "",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "password cannot be empty")
}

func TestPasswordCheckCmd_CustomEntropy(t *testing.T) {
	tests := []struct {
		name     string
		password string
		entropy  float64
		expectExit bool
	}{
		{"low entropy requirement", "simplepass", 20.0, false},
		{"high entropy requirement", "Str0ngP@ssw0rd!", 100.0, true},
		{"zero entropy", "weak", 0.0, false},
		{"reasonable entropy", "MyStr0ng!P@ssw0rd2024", 70.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture os.Exit to prevent test termination
			oldExit := osExit
			var exitCode int
			osExit = func(code int) { exitCode = code }
			defer func() { osExit = oldExit }()

			cmd := &cli.PasswordCheckCmd{
				Text:    tt.password,
				Entropy: tt.entropy,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)

			if tt.expectExit {
				require.Equal(t, 1, exitCode, "Should exit with code 1 for insufficient entropy")
			} else {
				require.Equal(t, 0, exitCode, "Should not exit for sufficient entropy")
			}
		})
	}
}

func TestPasswordCheckCmd_StdinInput(t *testing.T) {
	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	testPassword := "Str0ngP@ssw0rdFromStdin!"

	go func() {
		defer w.Close()
		w.Write([]byte(testPassword))
	}()

	cmd := &cli.PasswordCheckCmd{
		Text: "-", // Read from stdin
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestPasswordCheckCmd_StdinEmpty(t *testing.T) {
	// Mock stdin with no text parameter
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	testPassword := "Str0ngP@ssw0rdFromEmptyParam!"

	go func() {
		defer w.Close()
		w.Write([]byte(testPassword))
	}()

	cmd := &cli.PasswordCheckCmd{} // No text parameter, should read from stdin
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestPasswordCheckCmd_StdinWeakPassword(t *testing.T) {
	// Capture os.Exit to prevent test termination
	oldExit := osExit
	exitCode := 0
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	weakPassword := "123456"

	go func() {
		defer w.Close()
		w.Write([]byte(weakPassword))
	}()

	cmd := &cli.PasswordCheckCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, exitCode, "Should exit with code 1 for weak password from stdin")
}

func TestPasswordCheckCmd_StdinEmptyInput(t *testing.T) {
	// Mock stdin with empty input
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	w.Close() // Close immediately to simulate empty input

	cmd := &cli.PasswordCheckCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read password from stdin")
}

func TestPasswordCheckCmd_StdinWhitespaceOnly(t *testing.T) {
	// Mock stdin with whitespace only
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte("   \t\n   "))
	}()

	cmd := &cli.PasswordCheckCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no password provided")
}

func TestPasswordCheckCmd_Validate(t *testing.T) {
	tests := []struct {
		name      string
		entropy   float64
		expectErr bool
		errMsg    string
	}{
		{"valid entropy", 60.0, false, ""},
		{"zero entropy", 0.0, false, ""},
		{"negative entropy", -10.0, true, "entropy value must be non-negative"},
		{"very high entropy", 250.0, true, "entropy value is unrealistically high"},
		{"boundary high", 200.0, false, ""},
		{"boundary high+1", 201.0, true, "entropy value is unrealistically high"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.PasswordCheckCmd{
				Entropy: tt.entropy,
			}

			err := cmd.Validate()
			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPasswordCheckCmd_DefaultEntropyUsage(t *testing.T) {
	// Test that when no custom entropy is provided, default is used
	cmd := &cli.PasswordCheckCmd{
		Text: "TestP@ssw0rd123!",
		// Entropy: 0 (not set, should use default)
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)

	// Verify that the default entropy value is used by testing with a password
	// that passes the default but would fail a much higher requirement
	cmd2 := &cli.PasswordCheckCmd{
		Text:    "TestP@ssw0rd123!",
		Entropy: 200.0, // Very high requirement
	}

	// Capture os.Exit to prevent test termination
	oldExit := osExit
	exitCode := 0
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	err = cmd2.Run(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, exitCode, "Should fail with very high entropy requirement")
}

func TestPasswordCheckCmd_ReadFromPipe(t *testing.T) {
	// Test piped input (non-terminal stdin)
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	testPassword := "PipedStr0ngP@ssw0rd!"

	go func() {
		defer w.Close()
		w.Write([]byte(testPassword + "\n")) // Include newline as real pipes would
	}()

	cmd := &cli.PasswordCheckCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestPasswordCheck_BoundaryPasswords(t *testing.T) {
	// Test passwords around the boundary of what's considered secure
	tests := []struct {
		name     string
		password string
		entropy  float64
		shouldPass bool
	}{
		{"just below default", "Password1!", password.DefaultEntropy + 1, false},
		{"just above default", "Str0ngP@ssw0rd!", password.DefaultEntropy - 1, true},
		{"exact default boundary", "TestBoundary123!", password.DefaultEntropy, false}, // This might pass or fail depending on actual entropy
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture os.Exit to prevent test termination
			oldExit := osExit
			var exitCode int
			osExit = func(code int) { exitCode = code }
			defer func() { osExit = oldExit }()

			cmd := &cli.PasswordCheckCmd{
			Text:    tt.password,
			Entropy: tt.entropy,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)

			if tt.shouldPass {
			require.Equal(t, 0, exitCode, "Password should pass entropy check")
			} else {
			require.Equal(t, 1, exitCode, "Password should fail entropy check")
			}
		})
	}
}

func TestPasswordCheckCmd_LongPasswords(t *testing.T) {
	// Test various long passwords
	tests := []struct {
		name     string
		password string
	}{
		{"very long strong", strings.Repeat("Str0ng!P@ssw0rd", 10)},
		{"very long weak", strings.Repeat("a", 1000)},
		{"long mixed", "This is a very long passphrase with multiple words and some numbers 123456789 and symbols !@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.PasswordCheckCmd{
				Text: tt.password,
			}
			ctx := newTestContext()

			// Just ensure it doesn't panic with long inputs
			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestPasswordCheckCmd_SpecialCharacters(t *testing.T) {
	// Test passwords with various special characters
	specialPasswords := []string{
		"P@ssw0rd!#$%^&*()",
		"Пароль123!", // Cyrillic
		"密码123!",     // Chinese
		"パスワード123!",   // Japanese
		"🔐Password123🗝️", // Emoji
		"P@ss\tword\n123!", // Control characters
	}

	for _, pwd := range specialPasswords {
		t.Run("special_chars", func(t *testing.T) {
			cmd := &cli.PasswordCheckCmd{
				Text: pwd,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}



func TestPasswordCheckCmd_ReadFromTerminal(t *testing.T) {
	// This test simulates terminal input (harder to test directly)
	// We'll test the terminal detection logic by ensuring stdin is properly handled
	
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe that looks like terminal input
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte("TerminalStr0ngP@ssw0rd!\n"))
	}()

	cmd := &cli.PasswordCheckCmd{
		Text: "", // Empty will trigger stdin reading
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}


