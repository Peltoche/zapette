package systats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/clock"
)

type SystatsCron struct {
	service Service
	clock   clock.Clock
}

func newSystatCron(service Service, tools tools.Tools) *SystatsCron {
	return &SystatsCron{
		service: service,
		clock:   tools.Clock(),
	}
}

func (c *SystatsCron) Name() string {
	return "systats"
}

func (c *SystatsCron) Duration() time.Duration {
	return 300 * time.Millisecond
}

func (c *SystatsCron) Run(ctx context.Context) error {
	now := c.clock.Now()

	if now.Second()%5 != 0 {
		// The cron run every 5 seconds based on the computer clock.
		return nil
	}

	latest, err := c.service.GetLatest(ctx)
	if err != nil && !errors.Is(err, errNotFound) {
		return fmt.Errorf("failed to get the latest stats: %w", err)
	}

	if latest != nil && now.Sub(latest.time).Seconds() < 5.0 {
		// This 5s period have been done already.
		return nil
	}

	_, err = c.service.fetchAndRegister(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch the stats: %w", err)
	}

	return nil
}
