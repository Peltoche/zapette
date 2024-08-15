package sysstats

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
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
	GetRange(ctx context.Context, start time.Time, end time.Time) ([]Stats, error)
}

type service struct {
	fs          afero.Fs
	storage     storage
	clock       clock.Clock
	watchers    []chan struct{}
	watcherLock *sync.Mutex
}

func newService(storage storage, fs afero.Fs, tools tools.Tools) *service {
	return &service{
		storage:     storage,
		fs:          fs,
		clock:       tools.Clock(),
		watchers:    []chan struct{}{},
		watcherLock: new(sync.Mutex),
	}
}

func (s *service) SQLHookName() string {
	return "sysstats-svc"
}

// RunHook run as a hook for any db update (insert, update, delete).
func (s *service) RunSQLHook(ctx context.Context, table string) error {
	if table != "sysstats" {
		return nil
	}

	s.watcherLock.Lock()
	defer s.watcherLock.Unlock()

	// Try to fill the chan and skip the event if an one already needs to before
	// processed.
	for _, watcher := range s.watchers {
		select {
		case watcher <- struct{}{}:
		default:
		}
	}

	return nil
}

func (s *service) Watch(ctx context.Context) chan struct{} {
	c := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()
		close(c)

		s.watcherLock.Lock()
		defer s.watcherLock.Unlock()
		s.watchers = slices.DeleteFunc[[]chan struct{}](s.watchers, func(n chan struct{}) bool {
			return n == c
		})
	}()

	s.watcherLock.Lock()
	defer s.watcherLock.Unlock()
	s.watchers = append(s.watchers, c)

	return c
}

func (s *service) GetLast5mn(ctx context.Context) ([]Stats, error) {
	now := s.clock.Now()
	start := now.Add(-5 * time.Minute)

	stats, err := s.storage.GetRange(ctx, start, now)
	if err != nil {
		return nil, err
	}

	res := make([]Stats, (5*time.Minute)/(5*time.Second))

	for _, stat := range stats {
		res[int(stat.Time().Sub(start).Seconds())/5] = stat
	}

	return res, nil
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

func (s *service) fetch(_ context.Context) (*Stats, error) {
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
