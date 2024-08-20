package sysstats

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystats(t *testing.T) {
	t.Run("Marshal/Unmarshal binary", func(t *testing.T) {
		stats := NewFakeStats(t).Build()

		var buf []byte

		t.Run("MarshalBinary success", func(t *testing.T) {
			var err error

			buf, err = stats.MarshalBinary()
			require.NoError(t, err)
		})

		t.Run("UnmarshalBinary success", func(t *testing.T) {
			res := &Stats{}

			err := res.UnmarshalBinary(buf)
			require.NoError(t, err)

			assert.EqualValues(t, stats, res)
		})
	})

	t.Run("MarshalJSON success", func(t *testing.T) {
		stats := NewFakeStats(t).Build()

		buf, err := stats.MarshalJSON()
		require.NoError(t, err)

		assert.JSONEq(t, fmt.Sprintf(`{
			"time": "%s",
			"memory": {
				"totalMem": %.2f,
				"totalSwap": %.2f,
				"availableMem": %.2f,
				"buffers": %.2f,
				"cached": %.2f,
				"freeMem": %.2f,
				"freeSwap": %.2f,
				"shmem": %.2f,
				"sReclaimable": %.2f
			}
		}`,
			stats.time.Format(time.RFC3339),
			stats.memory.totalMem.GBytes(),
			stats.memory.totalSwap.GBytes(),
			stats.memory.availableMem.GBytes(),
			stats.memory.buffers.GBytes(),
			stats.memory.cached.GBytes(),
			stats.memory.freeMem.GBytes(),
			stats.memory.freeSwap.GBytes(),
			stats.memory.shmem.GBytes(),
			stats.memory.sReclaimable.GBytes(),
		), string(buf))
	})
}
