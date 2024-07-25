package systats

import (
	"testing"

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
}
