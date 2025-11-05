package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code"
)

func main() {
	cmd := cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "hexlet-path-size - print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		UsageText: "hexlet-path-size [global options]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("path is required")
			}
			path := c.Args().Get(0)

			allFiles := c.Bool("all")
			recursive := c.Bool("recursive")
			human := c.Bool("human")
			pathSize, err := code.GetPathSize(path, recursive, human, allFiles)
			if err != nil {
				return err
			}
			fmt.Printf("%s\t%s\n", pathSize, path)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "human",
				Aliases:     []string{"H"},
				DefaultText: "false",
				Usage:       "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				DefaultText: "false",
				Usage:       "include hidden files and directories",
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				DefaultText: "false",
				Usage:       "recursive size of directories",
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		return
	}
}
