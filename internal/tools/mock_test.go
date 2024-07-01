package tools

import (
	"log/slog"
	"testing"

	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/Peltoche/zapette/internal/tools/password"
	"github.com/Peltoche/zapette/internal/tools/response"
	"github.com/Peltoche/zapette/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMockToolbox(t *testing.T) {
	tools := NewMock(t)

	assert.IsType(t, new(clock.Mock), tools.Clock())
	assert.IsType(t, new(uuid.Mock), tools.UUID())
	assert.IsType(t, new(response.Mock), tools.ResWriter())

	assert.IsType(t, new(slog.Logger), tools.Logger())
	assert.IsType(t, new(password.Mock), tools.Password())
}
