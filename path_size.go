package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := GetSize(path, all, recursive)
	if err != nil {
		return "", err
	}

	return FormatSize(size, human) + "\t" + path, nil
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
