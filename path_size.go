package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := getSize(path, all, recursive)
	if err != nil {
		return "", err
	}

	return formatSize(size, human), nil
}

func getSize(path string, allFiles bool, recursive bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			return 0, err
		}
		return getSize(target, allFiles, recursive)
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
		size, err := entrySize(dir, entry, allFiles, recursive)
		if err != nil {
			return 0, err
		}
		total += size
	}

	return total, nil
}

func entrySize(parent string, entry os.DirEntry, allFiles bool, recursive bool) (int64, error) {
	name := entry.Name()

	if !allFiles && strings.HasPrefix(name, ".") {
		return 0, nil
	}

	fullPath := filepath.Join(parent, name)

	info, err := entry.Info()
	if err != nil {
		return 0, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return symlinkSize(fullPath, allFiles, recursive)
	}
	if info.Mode().IsRegular() {
		return info.Size(), nil
	}
	if info.IsDir() {
		if !recursive {
			return 0, nil
		}
		return dirSize(fullPath, allFiles, recursive)
	}
	return 0, nil
}

func symlinkSize(path string, allFiles bool, recursive bool) (int64, error) {
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return 0, err
	}
	return getSize(target, allFiles, recursive)
}

func formatSize(size int64, human bool) string {
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
