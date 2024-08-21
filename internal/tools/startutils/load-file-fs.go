package startutils

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func LoadFileinFS(t *testing.T, destFS afero.Fs, sourcePath, destPath string) {
	t.Helper()

	rawFile, err := os.ReadFile(sourcePath)
	require.NoError(t, err)

	err = afero.WriteFile(destFS, destPath, rawFile, 0o644)
	require.NoError(t, err)
}
