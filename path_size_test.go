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
