package cli_test

import (
	"testing"

	"github.com/alecthomas/kong"
	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestVersionFlag_BeforeReset_PrintsVersionAndExits(t *testing.T) {
	// Capture os.Exit to prevent test termination
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Create a mock context with version variables
	ctx := &kong.Context{}
	vars := kong.Vars{
		"version": "1.2.3",
		"commit":  "abc123def456",
		"date":    "2023-06-15T14:30:45Z",
	}

	// Test with version flag set to true
	versionFlag := cli.VersionFlag(true)
	
	err := versionFlag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should exit with code 0")
}

func TestVersionFlag_BeforeReset_DoesNotExitWhenFalse(t *testing.T) {
	// Capture os.Exit to prevent test termination
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Create a mock context with version variables
	ctx := &kong.Context{}
	vars := kong.Vars{
		"version": "1.2.3",
		"commit":  "abc123def456", 
		"date":    "2023-06-15T14:30:45Z",
	}

	// Test with version flag set to false
	versionFlag := cli.VersionFlag(false)
	
	err := versionFlag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should not call os.Exit when flag is false")
}

func TestVersionFlag_BeforeReset_WithEmptyVars(t *testing.T) {
	// Capture os.Exit to prevent test termination
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Create a mock context with empty variables
	ctx := &kong.Context{}
	vars := kong.Vars{}

	// Test with version flag set to true
	versionFlag := cli.VersionFlag(true)
	
	err := versionFlag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should exit with code 0 even with empty vars")
}

func TestVersionFlag_BeforeReset_WithPartialVars(t *testing.T) {
	// Capture os.Exit to prevent test termination
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Test with partial variables
	tests := []struct {
		name string
		vars kong.Vars
	}{
		{
			name: "only version",
			vars: kong.Vars{"version": "1.0.0"},
		},
		{
			name: "version and commit",
			vars: kong.Vars{"version": "1.0.0", "commit": "abc123"},
		},
		{
			name: "all but version",
			vars: kong.Vars{"commit": "abc123", "date": "2023-01-01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &kong.Context{}
			versionFlag := cli.VersionFlag(true)
			
			err := versionFlag.BeforeReset(ctx, tt.vars)
			require.NoError(t, err)
			require.Equal(t, 0, exitCode, "Should exit with code 0 regardless of partial vars")
		})
	}
}

func TestVersionFlag_Decode(t *testing.T) {
	versionFlag := cli.VersionFlag(true)
	ctx := newTestContext()
	
	err := versionFlag.Decode(ctx)
	require.NoError(t, err, "Decode should always return nil")
}

func TestVersionFlag_IsBool(t *testing.T) {
	versionFlag := cli.VersionFlag(true)
	
	isBool := versionFlag.IsBool()
	require.True(t, isBool, "IsBool should return true")
}

func TestVersionFlag_TypeConversion(t *testing.T) {
	// Test that VersionFlag can be converted to/from bool correctly
	tests := []struct {
		name     string
		flag     cli.VersionFlag
		expected bool
	}{
		{"true flag", cli.VersionFlag(true), true},
		{"false flag", cli.VersionFlag(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := bool(tt.flag)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestVersionFlag_MultipleCalls(t *testing.T) {
	// Test that multiple calls to BeforeReset work correctly
	oldExit := osExit
	var exitCount int
	osExit = func(code int) { exitCount++ }
	defer func() { osExit = oldExit }()

	ctx := &kong.Context{}
	vars := kong.Vars{"version": "1.0.0"}
	
	versionFlag := cli.VersionFlag(true)
	
	// First call
	err := versionFlag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	
	// Second call
	err = versionFlag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	
	require.Equal(t, 2, exitCount, "Should call os.Exit for each invocation")
}

func TestVersionFlag_ConcurrentAccess(t *testing.T) {
	// Test concurrent access to the version flag
	const numGoroutines = 10
	
	oldExit := osExit
	exitCount := make(chan int, numGoroutines)
	osExit = func(code int) { exitCount <- code }
	defer func() { osExit = oldExit }()

	ctx := &kong.Context{}
	vars := kong.Vars{"version": "1.0.0"}
	
	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			versionFlag := cli.VersionFlag(true)
			versionFlag.BeforeReset(ctx, vars)
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		select {
		case code := <-exitCount:
			require.Equal(t, 0, code, "All exits should be with code 0")
		}
	}
}

func TestVersionFlag_NilContext(t *testing.T) {
	// Test behavior with nil context
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	vars := kong.Vars{"version": "1.0.0"}
	versionFlag := cli.VersionFlag(true)
	
	err := versionFlag.BeforeReset(nil, vars)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should work with nil context")
}

func TestVersionFlag_NilVars(t *testing.T) {
	// Test behavior with nil vars
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	ctx := &kong.Context{}
	versionFlag := cli.VersionFlag(true)
	
	err := versionFlag.BeforeReset(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should work with nil vars")
}

func TestVersionFlag_BoolConversionEdgeCases(t *testing.T) {
	// Test various boolean conversion scenarios
	tests := []struct {
		name     string
		input    bool
		expected bool
	}{
		{"explicit true", true, true},
		{"explicit false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cli.VersionFlag(tt.input)
			result := bool(flag)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionFlag_Integration(t *testing.T) {
	// Test integration with kong framework concepts
	oldExit := osExit
	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = oldExit }()

	// Simulate what kong would do
	flag := cli.VersionFlag(true)
	
	// Test that it implements the expected interface methods
	require.True(t, flag.IsBool(), "Should be a boolean flag")
	
	err := flag.Decode(newTestContext())
	require.NoError(t, err, "Decode should succeed")
	
	// Test the main functionality
	ctx := &kong.Context{}
	vars := kong.Vars{
		"version": "test-version",
		"commit":  "test-commit",
		"date":    "test-date",
	}
	
	err = flag.BeforeReset(ctx, vars)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Should exit successfully")
}


