package code_test

import (
	"code"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize_Bytes(t *testing.T) {

	filePath := filepath.Join("testdata/fixture/dir1", "test.txt")

	data := make([]byte, 2048)

	err := os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)

	size, err := code.GetPathSize(filePath, false, false, false)
	require.NoError(t, err)

	require.Equal(t, "2048B", size)
}

func TestGetPathSize_HumanReadableKilobytes(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "test.txt")

	data := make([]byte, 2048)

	err := os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)

	size, err := code.GetPathSize(filePath, true, false, false)
	require.NoError(t, err)

	require.Equal(t, "2.0KB", size)
}

func TestGetPathSize_Recursive(t *testing.T) {
	tmpDir := t.TempDir()

	// файл первого уровня: 2 байта
	err := os.WriteFile(
		filepath.Join(tmpDir, "a.txt"),
		[]byte("ab"),
		0644,
	)
	require.NoError(t, err)

	// вложенная директория
	nestedDir := filepath.Join(tmpDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	require.NoError(t, err)

	// файл во вложенной директории: 5 байт
	err = os.WriteFile(
		filepath.Join(nestedDir, "b.txt"),
		[]byte("abcde"),
		0644,
	)
	require.NoError(t, err)

	size, err := code.GetPathSize(tmpDir, true, false, true)
	require.NoError(t, err)

	require.Equal(t, "7.0B", size)
}
