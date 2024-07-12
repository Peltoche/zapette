package systats

import (
	"context"
	"os"
	"testing"

	"github.com/Peltoche/zapette/internal/tools/datasize"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchMemInfos(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		afs := afero.NewMemMapFs()
		loadFileinFS(t, afs, "./testdata/meminfo.txt", "/proc/meminfo")

		svc := newService(afs)

		res, err := svc.FetchMeminfos(context.Background())
		require.NoError(t, err)
		assert.Equal(t, res, &Memory{
			totalMem:     datasize.ByteSize(16199860224),
			availableMem: datasize.ByteSize(13072998400),
			freeMem:      datasize.ByteSize(11920449536),
			buffers:      datasize.ByteSize(127262720),
			cached:       datasize.ByteSize(1712652288),
			sReclaimable: datasize.ByteSize(222531584),
			shmem:        datasize.ByteSize(557780992),
			totalSwap:    datasize.ByteSize(4294963200),
			freeSwap:     datasize.ByteSize(4294963200),
		})

		assert.Equal(t, "15.1 GB", res.TotalMemory().HumanReadable())
		assert.Equal(t, "11.1 GB", res.FreeMemory().HumanReadable())
		assert.Equal(t, "2.9 GB", res.UsedMemory().HumanReadable())
	})
}

func loadFileinFS(t *testing.T, destFS afero.Fs, sourcePath, destPath string) {
	t.Helper()

	rawFile, err := os.ReadFile(sourcePath)
	require.NoError(t, err)

	err = afero.WriteFile(destFS, destPath, rawFile, 0o644)
	require.NoError(t, err)
}
