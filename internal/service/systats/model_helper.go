package systats

import (
	"testing"
	"time"

	"github.com/Peltoche/zapette/internal/tools/datasize"
	"github.com/brianvoe/gofakeit/v7"
)

type FakeStatsBuilder struct {
	t     testing.TB
	stats *Stats
}

func NewFakeStats(t testing.TB) *FakeStatsBuilder {
	t.Helper()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	totalMem := datasize.ByteSize(gofakeit.Number(int(datasize.GB), 20*int(datasize.GB)))
	totalSwap := datasize.ByteSize(gofakeit.Number(int(datasize.GB), 20*int(datasize.GB)))

	return &FakeStatsBuilder{
		t: t,
		stats: &Stats{
			time: createdAt,
			memory: &Memory{
				totalMem:     totalMem,
				availableMem: totalMem / 100 * 90,
				freeMem:      datasize.ByteSize(gofakeit.Number(0, int(totalMem))),
				buffers:      datasize.MB * 100,
				cached:       datasize.MB * 200,
				sReclaimable: datasize.MB * 20,
				shmem:        datasize.MB * 10,
				totalSwap:    totalSwap,
				freeSwap:     datasize.ByteSize(gofakeit.Number(0, int(totalSwap))),
			},
		},
	}
}

func (b *FakeStatsBuilder) WithTime(t time.Time) *FakeStatsBuilder {
	b.stats.time = t

	return b
}

func (b *FakeStatsBuilder) Build() *Stats {
	return b.stats
}
