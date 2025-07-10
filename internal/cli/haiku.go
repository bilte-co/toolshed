package cli

import (
	"fmt"

	"github.com/bilte-co/toolshed/haiku"
)

type HaikuCmd struct {
	Generate HaikuGenerateCmd `cmd:"generate" help:"Generate a random haiku name"`
}

type HaikuGenerateCmd struct {
	Token int64  `short:"t" default:"9999" help:"Maximum token value (default: 9999)"`
	Delim string `short:"d" default:"-" help:"Delimiter between words (default: '-')"`
}

func (cmd *HaikuGenerateCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating haiku name")

	haikuinator := haiku.NewHaikunator()

	haikuName := haikuinator.TokenDelimHaikunate(cmd.Token, cmd.Delim)

	fmt.Println(haikuName)
	ctx.Logger.Info("Haiku name generated successfully", "haiku", haikuName)
	return nil
}
