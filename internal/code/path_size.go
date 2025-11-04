package code

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func BuildApp() {
	cmd := cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory;",
		UsageText: "hexlet-path-size [global options]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("path is required")
			}
			path := c.Args().Get(0)

			size, err := GetSize(path)
			if err != nil {
				return err
			}
			fmt.Printf("%d\t%s\n", size, path)
			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		return
	}
}

func GetSize(path string) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			fi, err := entry.Info()
			if err != nil {
				return 0, err
			}
			total += fi.Size()
		}
	}

	return total, nil
}
