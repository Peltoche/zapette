package sysinfos

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchSysInfos(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		toolsMock := tools.NewMock(t)
		afs := afero.NewMemMapFs()
		loadFileinFS(t, afs, "./testdata/uptime.txt", "/proc/uptime")
		loadFileinFS(t, afs, "./testdata/hostname.txt", "/etc/hostname")

		svc := newService(afs, toolsMock)

		err := svc.fetch(context.Background())
		require.NoError(t, err)

		assert.Equal(t, "zapettePC", svc.hostname)
		expected, err := time.ParseDuration("109118.83s")
		require.NoError(t, err)
		assert.Equal(t, expected, svc.uptime)
	})
}

func loadFileinFS(t *testing.T, destFS afero.Fs, sourcePath, destPath string) {
	t.Helper()

	rawFile, err := os.ReadFile(sourcePath)
	require.NoError(t, err)

	err = afero.WriteFile(destFS, destPath, rawFile, 0o644)
	require.NoError(t, err)
}
