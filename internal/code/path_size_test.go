package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize_File(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "alpha.txt")

	got, err := GetSize(p)
	require.NoError(t, err)

	const want int64 = 3
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Empty(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "empty.txt")

	got, err := GetSize(p)
	require.NoError(t, err)

	const want int64 = 0
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Error(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "unknown.txt")

	_, err := GetSize(p)
	require.Error(t, err)
}

func TestGetPathSize_Directory(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "directory")

	got, err := GetSize(p)
	require.NoError(t, err)

	const want int64 = 6
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}
