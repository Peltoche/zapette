package tools

import (
	"log/slog"
	"testing"

	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/Peltoche/zapette/internal/tools/logger"
	"github.com/Peltoche/zapette/internal/tools/password"
	"github.com/Peltoche/zapette/internal/tools/response"
	"github.com/Peltoche/zapette/internal/tools/uuid"
	"github.com/neilotoole/slogt"
)

type Config struct {
	Log      logger.Config   `json:"log"`
	Response response.Config `json:"response"`
}

type Toolbox struct {
	clock     clock.Clock
	uuid      uuid.Service
	resWriter response.Writer
	log       *slog.Logger
	password  password.Password
}

func NewToolbox(cfg Config) *Toolbox {
	log := logger.NewSLogger(cfg.Log)

	return &Toolbox{
		clock:     clock.NewDefault(),
		uuid:      uuid.NewProvider(),
		log:       log,
		resWriter: response.Init(cfg.Response),
		password:  password.NewArgon2IDPassword(),
	}
}

func NewToolboxForTest(t testing.TB) *Toolbox {
	t.Helper()

	log := slogt.New(t)
	return &Toolbox{
		clock:     clock.NewDefault(),
		uuid:      uuid.NewProvider(),
		log:       log,
		resWriter: response.Init(response.Config{PrettyRender: true}),
		password:  password.NewArgon2IDPassword(),
	}
}

// Clock implements App.
//
// Return a clock.Default.
func (d *Toolbox) Clock() clock.Clock {
	return d.clock
}

// UUID implements App.
//
// Return a *uuid.Default.
func (d *Toolbox) UUID() uuid.Service {
	return d.uuid
}

// Logger implements App.
//
// Return a *logging.StdLogger.
func (d *Toolbox) Logger() *slog.Logger {
	return d.log
}

// ResWriter implements App.
//
// Return a *response.Writer.
func (d *Toolbox) ResWriter() response.Writer {
	return d.resWriter
}

func (d *Toolbox) Password() password.Password {
	return d.password
}
