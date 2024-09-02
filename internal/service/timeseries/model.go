package timeseries

import (
	"time"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation"
)

type Timeserie struct {
	tsID      uuid.UUID
	graphSpan time.Duration
	tickSpan  time.Duration
}

type TimeData struct {
	tsID uuid.UUID
	at   time.Time
	data []byte
}

type CreateCmd struct {
	CreatedBy *users.User
	GraphSpan time.Duration
	TickSpan  time.Duration
}

func (t CreateCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.CreatedBy, v.Required),
		v.Field(&t.GraphSpan, v.Required),
		v.Field(&t.TickSpan, v.Required),
	)
}
