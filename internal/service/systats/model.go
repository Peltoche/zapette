package systats

import (
	"math"
	"time"

	"github.com/Peltoche/zapette/internal/tools/datasize"
)

type Stats struct {
	time   time.Time
	memory *Memory
}

func (s *Stats) Time() time.Time {
	return s.time
}

func (s *Stats) Memory() *Memory {
	return s.memory
}

// Memory holds information on system memory usage
type Memory struct {
	totalMem     datasize.ByteSize
	availableMem datasize.ByteSize
	freeMem      datasize.ByteSize
	buffers      datasize.ByteSize
	cached       datasize.ByteSize
	sReclaimable datasize.ByteSize
	shmem        datasize.ByteSize
	totalSwap    datasize.ByteSize
	freeSwap     datasize.ByteSize
}

func (c Memory) TotalMemory() datasize.ByteSize {
	return c.totalMem
}

func (c Memory) AvailableMemory() datasize.ByteSize {
	return c.availableMem
}

func (c Memory) FreeMemory() datasize.ByteSize {
	return c.freeMem
}

// Cached correspond to the yellow bars in htop.
func (c Memory) BufCache() datasize.ByteSize {
	return c.buffers + c.cached + c.sReclaimable
}

func (c Memory) UsedMemory() datasize.ByteSize {
	return c.totalMem - c.availableMem
}

func (c Memory) PercentageUsedMemory() int {
	return int(math.Round(float64(c.UsedMemory()) / float64(c.totalMem) * 100))
}

func (c Memory) PercentageFreeMemory() int {
	return int(math.Round(float64(c.freeMem) / float64(c.totalMem) * 100))
}

func (c Memory) PercentageAvailableMemory() int {
	return int(math.Round(float64(c.availableMem) / float64(c.totalMem) * 100))
}

func (c Memory) TotalSwap() datasize.ByteSize {
	return c.totalSwap
}

func (c Memory) UsedSwap() datasize.ByteSize {
	return c.totalSwap - c.freeSwap
}
