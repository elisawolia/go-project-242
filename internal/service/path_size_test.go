package service

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize_File(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "alpha.txt")

	got, err := GetSize(p, true, false)
	require.NoError(t, err)

	const want int64 = 3
	require.Equalf(t, want, got, "got %v, want %v", got, want)

	got, err = GetSize(p, false, false)
	require.NoError(t, err)

	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Empty(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "empty.txt")

	got, err := GetSize(p, false, false)
	require.NoError(t, err)

	const want int64 = 0
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Error(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "unknown.txt")

	_, err := GetSize(p, false, false)
	require.Error(t, err)
}

func TestGetPathSize_Directory(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "directory")

	got, err := GetSize(p, false, false)
	require.NoError(t, err)

	const want int64 = 6
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Directory_AllFiles(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "directory")

	got, err := GetSize(p, true, false)
	require.NoError(t, err)

	const want int64 = 16
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetPathSize_Directory_Recursive(t *testing.T) {
	t.Parallel()

	p := filepath.Join("..", "..", "testdata", "directory")

	got, err := GetSize(p, false, true)
	require.NoError(t, err)

	const want int64 = 12
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestFormatSize_NotHuman(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		{"zero", 0, "0B"},
		{"bytes", 123, "123B"},
		{"kb_exact", 1024, "1024B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSize(tt.in, false)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormatSize_Human(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int64
		want string
	}{
		{"bytes", 123, "123B"},
		{"one_kb", 1024, "1.0KB"},
		{"half_mb", 512 * 1024, "512.0KB"},
		{"one_mb", 1024 * 1024, "1.0MB"},
		{"example_from_task", 1234567, "1.2MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSize(tt.in, true)
			require.Equal(t, tt.want, got)
		})
	}
}
