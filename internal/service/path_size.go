package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func BuildApp() {
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
			size, err := GetSize(path, allFiles, recursive)
			if err != nil {
				return err
			}
			human := c.Bool("human")
			fmt.Printf("%s\t%s\n", FormatSize(size, human), path)
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

func GetSize(path string, allFiles bool, recursive bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return info.Size(), nil
	}

	return dirSize(path, allFiles, recursive)
}

func dirSize(dir string, allFiles bool, recursive bool) (int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var total int64

	for _, entry := range entries {
		name := entry.Name()

		if !allFiles && strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dir, name)

		if entry.Type().IsRegular() {
			fi, err := entry.Info()
			if err != nil {
				return 0, err
			}
			total += fi.Size()
			continue
		}

		if entry.IsDir() {
			if recursive {
				subTotal, err := dirSize(fullPath, allFiles, recursive)
				if err != nil {
					return 0, err
				}
				total += subTotal
			}
			continue
		}
	}

	return total, nil
}

func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	val := float64(size)
	i := 0

	for val >= 1024 && i < len(units)-1 {
		val /= 1024
		i++
	}

	if i == 0 {
		return fmt.Sprintf("%d%s", size, units[i])
	}

	return fmt.Sprintf("%.1f%s", val, units[i])
}
