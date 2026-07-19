package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/dankobg/juicer/cmd/identities"
	"github.com/dankobg/juicer/cmd/juicer"
)

var CLI struct {
	Serve      juicer.ServeCommand `cmd:"" help:"Run Juicer server"`
	Identities identities.RootCmd  `cmd:"" help:"Manage identities"`
}

func Run() {
	c := kong.Parse(
		&CLI,
		kong.Name("juicer"),
		kong.Description("juicer chess server"),
	)

	err := c.Run()
	c.FatalIfErrorf(err)
}
