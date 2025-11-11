package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSize_File(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "alpha.txt")

	got, err := getSize(p, true, false)
	require.NoError(t, err)

	const want int64 = 3
	require.Equalf(t, want, got, "got %v, want %v", got, want)

	got, err = getSize(p, false, false)
	require.NoError(t, err)

	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetSize_Empty(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "empty.txt")

	got, err := getSize(p, false, false)
	require.NoError(t, err)

	const want int64 = 0
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetSize_Error(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "unknown.txt")

	_, err := getSize(p, false, false)
	require.Error(t, err)
}

func TestGetSize_Directory(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "directory")

	got, err := getSize(p, false, false)
	require.NoError(t, err)

	const want int64 = 6
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetSize_Directory_AllFiles(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "directory")

	got, err := getSize(p, true, false)
	require.NoError(t, err)

	const want int64 = 16
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetSize_Directory_Recursive(t *testing.T) {
	t.Parallel()

	p := filepath.Join("testdata", "directory")

	got, err := getSize(p, false, true)
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
			got := formatSize(tt.in, false)
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
			got := formatSize(tt.in, true)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetSize_SymlinkToFile(t *testing.T) {
	t.Parallel()

	target, err := filepath.Abs(filepath.Join("testdata", "alpha.txt"))
	require.NoError(t, err)

	tmp := t.TempDir()
	link := filepath.Join(tmp, "alpha-link.txt")

	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	got, err := getSize(link, false, false)
	require.NoError(t, err)

	const want int64 = 3
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestGetSize_SymlinkToDirectory(t *testing.T) {
	t.Parallel()

	target, err := filepath.Abs(filepath.Join("testdata", "directory"))
	require.NoError(t, err)

	tmp := t.TempDir()
	link := filepath.Join(tmp, "directory-link")

	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	got, err := getSize(link, false, true)
	require.NoError(t, err)

	const want int64 = 12
	require.Equalf(t, want, got, "got %v, want %v", got, want)
}

func TestEntrySize_SkipHidden(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	hidden := filepath.Join(tmp, ".secret")
	err := os.WriteFile(hidden, []byte("xxx"), 0o644)
	require.NoError(t, err)

	entry, err := os.ReadDir(tmp)
	require.NoError(t, err)
	require.Len(t, entry, 1)

	size, err := entrySize(tmp, entry[0], false, false)
	require.NoError(t, err)
	require.Equal(t, int64(0), size)

	size, err = entrySize(tmp, entry[0], true, false)
	require.NoError(t, err)
	require.Equal(t, int64(3), size)
}

func TestDirSize_NonRecursive(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	topFile := filepath.Join(tmp, "top.txt")
	err := os.WriteFile(topFile, []byte("abcd"), 0o644)
	require.NoError(t, err)

	sub := filepath.Join(tmp, "sub")
	err = os.Mkdir(sub, 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(sub, "inner.txt"), []byte("xyz"), 0o644)
	require.NoError(t, err)

	got, err := dirSize(tmp, false, false)
	require.NoError(t, err)
	require.Equal(t, int64(4), got)
}

func TestSymlinkSize(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	target := filepath.Join(tmp, "target.txt")
	err := os.WriteFile(target, []byte("hi"), 0o644)
	require.NoError(t, err)

	link := filepath.Join(tmp, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	got, err := symlinkSize(link, false, false)
	require.NoError(t, err)
	require.Equal(t, int64(2), got)
}
