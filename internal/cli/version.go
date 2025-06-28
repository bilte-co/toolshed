package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

// VersionFlag handles the --version flag
type VersionFlag bool

func (v VersionFlag) BeforeReset(ctx *kong.Context, vars kong.Vars) error {
	if bool(v) {
		fmt.Printf("toolshed %s\n", vars["version"])
		fmt.Printf("commit: %s\n", vars["commit"])
		fmt.Printf("built: %s\n", vars["date"])
		os.Exit(0)
	}
	return nil
}

func (v VersionFlag) Decode(ctx *CLIContext) error { return nil }
func (v VersionFlag) IsBool() bool                 { return true }
