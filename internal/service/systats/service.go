package systats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/Peltoche/zapette/internal/tools/datasize"
	"github.com/spf13/afero"
)

const filePath = "/proc/meminfo"

var (
	ErrInvalidFieldFormat = errors.New("invalid field format")
	ErrInvalidLineFormat  = errors.New("invalid line format")
	ErrUnsupportedUnit    = errors.New("unsupported unit")
)

func InvalidFieldFormat(key, expected, val string) error {
	return fmt.Errorf("%s: %w: expected an uint64, have %q", key, ErrInvalidFieldFormat, val)
}

type storage interface {
	GetLatest(ctx context.Context) (*Stats, error)
	Save(ctx context.Context, stats *Stats) error
}

type service struct {
	fs      afero.Fs
	storage storage
	clock   clock.Clock
}

func newService(storage storage, fs afero.Fs, tools tools.Tools) *service {
	return &service{
		storage: storage,
		fs:      fs,
		clock:   tools.Clock(),
	}
}

func (s *service) GetLatest(ctx context.Context) (*Stats, error) {
	return s.storage.GetLatest(ctx)
}

func (s *service) fetchAndRegister(ctx context.Context) (*Stats, error) {
	stats, err := s.fetch(ctx)
	if err != nil {
		return nil, err
	}

	err = s.storage.Save(ctx, stats)
	if err != nil {
		return nil, fmt.Errorf("failed to save the new stats: %w", err)
	}

	return stats, nil
}

func (s *service) fetch(ctx context.Context) (*Stats, error) {
	now := s.clock.Now().Truncate(time.Second)

	content, err := afero.ReadFile(s.fs, filePath)
	if err != nil {
		return nil, err
	}

	mem := Memory{}

	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			mem.totalMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "MemFree:":
			mem.freeMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "MemAvailable:":
			mem.availableMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Buffers:":
			mem.buffers, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Cached:":
			mem.cached, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SReclaimable:":
			mem.sReclaimable, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Shmem:":
			mem.shmem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SwapTotal:":
			mem.totalSwap, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SwapFree:":
			mem.freeSwap, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}
		}
	}

	stats := Stats{
		time:   now,
		memory: &mem,
	}

	return &stats, nil
}

func parseBytesValue(fields []string) (datasize.ByteSize, error) {
	if len(fields) != 3 {
		return 0, ErrInvalidLineFormat
	}

	res, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0, InvalidFieldFormat(fields[0], "uint64", fields[1])
	}

	switch fields[2] {
	case "kB":
		return datasize.ByteSize(float64(res) * 1024), nil
	case "mB":
		return datasize.ByteSize(float64(res) * 1024 * 1024), nil
	case "gB":
		return datasize.ByteSize(float64(res) * 1024 * 1024 * 1024), nil
	default:
		return 0, ErrUnsupportedUnit
	}
}
