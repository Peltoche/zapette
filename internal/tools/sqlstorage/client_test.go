package sqlstorage

import (
	"testing"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSQliteClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := Config{Path: t.TempDir() + "/db.sqlite"}
		hooks := []SQLChangeHook{}
		tools := tools.NewToolboxForTest(t)

		client, err := NewSQliteClient(&cfg, hooks, tools.Logger())
		require.NoError(t, err)

		require.NoError(t, client.Ping())
	})

	t.Run("with an invalid path", func(t *testing.T) {
		cfg := Config{Path: "/foo/some-invalidpath"}
		hooks := []SQLChangeHook{}
		tools := tools.NewToolboxForTest(t)

		client, err := NewSQliteClient(&cfg, hooks, tools.Logger())
		assert.Nil(t, client)
		require.EqualError(t, err, "unable to open database file: no such file or directory")
	})

	t.Run("with not specified path", func(t *testing.T) {
		cfg := Config{Path: ""}
		hooks := []SQLChangeHook{}
		tools := tools.NewToolboxForTest(t)

		client, err := NewSQliteClient(&cfg, hooks, tools.Logger())
		assert.NotNil(t, client)
		require.NoError(t, err)
	})
}
