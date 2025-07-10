package cli

import (
	"fmt"

	"github.com/bilte-co/toolshed/haiku"
)

type HaikuCmd struct {
	Generate HaikuGenerateCmd `cmd:"generate" help:"Generate a random haiku name"`
}

type HaikuGenerateCmd struct {
	Token   int64  `short:"t" default:"9999" help:"Maximum token value (default: 9999, use 0 for no token)"`
	Delim   string `short:"d" default:"-" help:"Delimiter between words (default: '-')"`
	NoToken bool   `short:"n" help:"Generate haiku without numeric token"`
}

func (cmd *HaikuGenerateCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating haiku name")

	haikuinator := haiku.NewHaikunator()

	var haikuName string
	var err error

	// Choose the appropriate method based on whether token is needed
	if cmd.NoToken || cmd.Token == 0 {
		haikuName, err = haikuinator.DelimHaikunate(cmd.Delim)
	} else {
		haikuName, err = haikuinator.TokenDelimHaikunate(cmd.Token, cmd.Delim)
	}

	if err != nil {
		ctx.Logger.Error("Failed to generate haiku name", "error", err)
		return fmt.Errorf("failed to generate haiku: %w", err)
	}

	fmt.Println(haikuName)
	ctx.Logger.Info("Haiku name generated successfully", "haiku", haikuName)
	return nil
}
