package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Peltoche/zapette/assets"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/logger"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/Peltoche/zapette/internal/tools/startutils"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

var testConfig = Config{
	FS:       afero.NewMemMapFs(),
	Listener: router.Config{},
	Assets:   assets.Config{},
	Storage:  sqlstorage.Config{Path: ":memory:"},
	Tools:    tools.Config{Log: logger.Config{Output: io.Discard}},
	HTML:     html.Config{},
	Folder:   "/foo",
}

func TestServerStart(t *testing.T) {
	ctx := context.Background()

	app := start(ctx, testConfig, fx.Invoke(func(*router.API) {}))
	require.NoError(t, app.Err())
}

func TestServerRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := startutils.GetFreePort(t)

	testConfig.Listener.Addr = fmt.Sprintf("localhost:%d", port)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var runErr error
	go func() {
		defer wg.Done()
		_, runErr = Run(ctx, testConfig)
	}()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/web/login", port), nil)
	require.NoError(t, err)

	var res *http.Response
	for range 50 {
		res, err = http.DefaultClient.Do(req)
		if err == nil || !strings.Contains(err.Error(), "connection refused") {
			break
		}

		if res != nil {
			res.Body.Close()
		}
		time.Sleep(20 * time.Millisecond)
	}

	cancel()
	wg.Wait()

	require.NoError(t, runErr)
}
