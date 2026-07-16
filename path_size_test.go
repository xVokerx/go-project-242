package code_test

import (
	"code"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize_Bytes(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "test.txt")

	data := make([]byte, 2048)

	err := os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)

	size, err := code.GetPathSize(filePath, false, false, false)
	require.NoError(t, err)

	require.Equal(t, "2048B", size)
}

func TestGetPathSize_Bytes_AllFlags(t *testing.T) {
	testDir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(testDir, "test.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	dir1 := filepath.Join(testDir, "dir1")
	err = os.Mkdir(dir1, 0755)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(dir1, ".test1.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	size, err := code.GetPathSize(testDir, true, true, true)
	require.NoError(t, err)

	require.Equal(t, "4.0KB", size)
}

func TestGetPathSize_NoExistFile(t *testing.T) {
	_, err := code.GetPathSize(
		"./test.txt",
		false,
		false,
		false,
	)

	require.Error(t, err)
}

func TestGetPathSize_NoPath(t *testing.T) {
	_, err := code.GetPathSize(
		"",
		false,
		false,
		false,
	)

	require.Error(t, err)
}

func TestGetPathSize_Bytes_Recursive(t *testing.T) {
	testDir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(testDir, "test.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	dir1 := filepath.Join(testDir, "dir1")
	err = os.Mkdir(dir1, 0755)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(dir1, "test1.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	size, err := code.GetPathSize(testDir, true, false, false)
	require.NoError(t, err)

	require.Equal(t, "4096B", size)

}

func TestGetPathSize_Bytes_All(t *testing.T) {
	testDir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(testDir, "test.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(testDir, ".test1.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	size, err := code.GetPathSize(testDir, false, false, true)
	require.NoError(t, err)

	require.Equal(t, "4096B", size)

}

func TestRecursiveDirSize(t *testing.T) {
	testDir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(testDir, "a.txt"),
		make([]byte, 100),
		0644,
	)
	require.NoError(t, err)

	subDir := filepath.Join(testDir, "sub")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(subDir, "b.txt"),
		make([]byte, 200),
		0644,
	)
	require.NoError(t, err)

	size, err := code.RecursiveDirSize(testDir, true)
	require.NoError(t, err)

	require.Equal(t, int64(300), size)
}

func TestGetPathSize_Bytes_IgnoreHidden(t *testing.T) {
	testDir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(testDir, "test.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(testDir, ".test1.txt"),
		make([]byte, 2048),
		0644,
	)
	require.NoError(t, err)

	size, err := code.GetPathSize(testDir, false, false, false)
	require.NoError(t, err)

	require.Equal(t, "2048B", size)

}

func TestHumanize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{
			name: "bytes",
			size: 7,
			want: "7B",
		},
		{
			name: "KBytes",
			size: 2048,
			want: "2.0KB",
		},
		{
			name: "MBytes",
			size: 1024 * 1024,
			want: "1.0MB",
		},
		{
			name: "GBytes",
			size: 1024 * 1024 * 1024,
			want: "1.0GB",
		},
		{
			name: "TBytes",
			size: 1024 * 1024 * 1024 * 1024,
			want: "1.0TB",
		},
		{
			name: "PBytes",
			size: 1024 * 1024 * 1024 * 1024 * 1024,
			want: "1.0PB",
		},
		{
			name: "EBytes",
			size: 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
			want: "1.0EB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, code.Humanize(tt.size))
		})
	}
}
