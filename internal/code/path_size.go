package code

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func BuildApp() {
	cmd := cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory;",
		UsageText: "hexlet-path-size [global options]",
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		return
	}
}
