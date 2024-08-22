package sysinfos

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/startutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchSysInfos(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		toolsMock := tools.NewMock(t)
		afs := afero.NewMemMapFs()
		startutils.LoadFileinFS(t, afs, "./testdata/uptime.txt", "/proc/uptime")
		startutils.LoadFileinFS(t, afs, "./testdata/hostname.txt", "/etc/hostname")

		now := time.Now()

		toolsMock.ClockMock.On("Now").Return(now).Once()

		svc := newService(afs, toolsMock)

		err := svc.fetch(context.Background())
		require.NoError(t, err)

		assert.Equal(t, "zapettePC", svc.hostname)

		expected := now.Add(-109118 * time.Second)
		assert.Equal(t, expected, svc.startTime)
	})
}
