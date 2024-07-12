package systats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

type service struct {
	fs afero.Fs
}

func newService(fs afero.Fs) *service {
	return &service{
		fs: fs,
	}
}

func (s *service) FetchMeminfos(ctx context.Context) (*Memory, error) {
	content, err := afero.ReadFile(s.fs, filePath)
	if err != nil {
		return nil, err
	}

	cpu := Memory{}

	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			cpu.totalMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "MemFree:":
			cpu.freeMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "MemAvailable:":
			cpu.availableMem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Buffers:":
			cpu.buffers, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Cached:":
			cpu.cached, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SReclaimable:":
			cpu.sReclaimable, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "Shmem:":
			cpu.shmem, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SwapTotal:":
			cpu.totalSwap, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}

		case "SwapFree:":
			cpu.freeSwap, err = parseBytesValue(fields)
			if err != nil {
				return nil, err
			}
		}
	}

	return &cpu, nil
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
