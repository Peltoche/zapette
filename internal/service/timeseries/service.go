package timeseries

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/errs"
	"github.com/Peltoche/zapette/internal/tools/uuid"
)

type storage interface {
	saveData(ctx context.Context, tsID uuid.UUID, t time.Time, data []byte) error
	getLatestData(ctx context.Context, tsID uuid.UUID) (*TimeData, error)

	saveConfig(ctx context.Context, cfg *Timeserie) error
	getAllConfigs(ctx context.Context) ([]Timeserie, error)
	getConfigByID(ctx context.Context, tsID uuid.UUID) (*Timeserie, error)
}

type service struct {
	storage storage
	uuid    uuid.Service
}

func newService(tools tools.Tools, storage storage) *service {
	return &service{
		storage: storage,
		uuid:    tools.UUID(),
	}
}

func (s *service) CreateTimeSerie(ctx context.Context, cmd *CreateCmd) (*Timeserie, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	if cmd.GraphSpan < time.Second {
		return nil, errs.Validation(fmt.Errorf("graphSpan must be at least 1s, have %s", cmd.GraphSpan))
	}

	if cmd.TickSpan < time.Second {
		return nil, errs.Validation(fmt.Errorf("tickSpan must be at least 1s, have %s", cmd.GraphSpan))
	}

	roundedGraphSpan := cmd.GraphSpan.Round(time.Second)
	roundedTickSpan := cmd.TickSpan.Round(time.Second)

	if cmd.TickSpan >= cmd.GraphSpan {
		return nil, errs.Validation(fmt.Errorf("the tick span must be shorter than the graph span, have %s >= %s", roundedTickSpan, roundedGraphSpan))
	}

	ts := Timeserie{
		tsID:      s.uuid.New(),
		graphSpan: roundedGraphSpan,
		tickSpan:  roundedTickSpan,
	}

	s.storage.saveConfig(ctx, &ts)

	return &ts, nil
}

func (s *service) GetTimeserieByID(ctx context.Context, id uuid.UUID) (*Timeserie, error) {
	return s.storage.getConfigByID(ctx, id)
}

func (s *service) GetAllTimeseries(ctx context.Context) ([]Timeserie, error) {
	return s.storage.getAllConfigs(ctx)
}

func (s *service) GetLatestData(ctx context.Context, ts *Timeserie) (*TimeData, error) {
	return s.storage.getLatestData(ctx, ts.tsID)
}

func (s *service) SaveLatestData(ctx context.Context, ts *Timeserie, at time.Time, data []byte) error {
	if at.Unix()%int64(ts.tickSpan.Seconds()) != 0 {
		return errors.New("the time doesn't correspond to the timeserie tick span")
	}

	latest, err := s.storage.getLatestData(ctx, ts.tsID)
	if err != nil && !errors.Is(err, errNotFound) {
		return fmt.Errorf("failed to tech the previous entry: %w", err)
	}

	if latest.at.After(at) || latest.at.Equal(at) {
		return errors.New("this tick have already been registered")
	}

	err = s.storage.saveData(ctx, ts.tsID, at, data)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	return nil
}
