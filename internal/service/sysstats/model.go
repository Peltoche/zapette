package sysstats

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/Peltoche/zapette/internal/tools/datasize"
)

type Namespace int

const (
	Unknown Namespace = iota
	MinGraph
)

type Stats struct {
	time   time.Time
	memory *Memory
}

func (s *Stats) Time() time.Time {
	return s.time.Truncate(time.Second)
}

func (s *Stats) Memory() *Memory {
	return s.memory
}

func (a *Stats) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.BigEndian, a.time.Truncate(time.Second).Unix())

	rawMemory, _ := a.memory.MarshalBinary()
	binary.Write(buf, binary.BigEndian, rawMemory)

	return buf.Bytes(), nil
}

func (a *Stats) IsEmpty() bool {
	return a.time.IsZero()
}

func (a *Stats) UnmarshalBinary(b []byte) error {
	unixSecs := int64(binary.BigEndian.Uint64(b))

	a.time = time.Unix(unixSecs, 0).UTC()

	a.memory = &Memory{}
	err := a.memory.UnmarshalBinary(b[8:])
	if err != nil {
		return fmt.Errorf("failed to decode the memory: %w", err)
	}

	return nil
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

func (c *Memory) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, c.totalMem)
	binary.Write(buf, binary.BigEndian, c.availableMem)
	binary.Write(buf, binary.BigEndian, c.freeMem)
	binary.Write(buf, binary.BigEndian, c.buffers)
	binary.Write(buf, binary.BigEndian, c.cached)
	binary.Write(buf, binary.BigEndian, c.sReclaimable)
	binary.Write(buf, binary.BigEndian, c.shmem)
	binary.Write(buf, binary.BigEndian, c.totalSwap)
	binary.Write(buf, binary.BigEndian, c.freeSwap)

	return buf.Bytes(), nil
}

func (c *Memory) UnmarshalBinary(b []byte) error {
	buf := bytes.NewReader(b)

	if err := binary.Read(buf, binary.BigEndian, &c.totalMem); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.availableMem); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.freeMem); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.buffers); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.cached); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.sReclaimable); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.shmem); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.totalSwap); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.freeSwap); err != nil {
		return err
	}

	return nil
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
