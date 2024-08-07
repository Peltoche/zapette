package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Peltoche/zapette/assets"
	"github.com/Peltoche/zapette/internal/migrations"
	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/service/utilities"
	"github.com/Peltoche/zapette/internal/service/websessions"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/cron"
	"github.com/Peltoche/zapette/internal/tools/logger"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	sysstatstmpl "github.com/Peltoche/zapette/internal/web/handlers/sysstats"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/middlewares"
	"github.com/spf13/afero"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type Folder string

type Config struct {
	fx.Out
	Tools    tools.Config
	FS       afero.Fs
	Storage  sqlstorage.Config
	Folder   Folder
	Listener router.Config
	HTML     html.Config
	Assets   assets.Config
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(router.Registerer)),
		fx.ResultTags(`group:"routes"`),
	)
}

func start(ctx context.Context, cfg Config, invoke fx.Option) *fx.App {
	app := fx.New(
		fx.WithLogger(func(tools tools.Tools) fxevent.Logger { return logger.NewFxLogger(tools.Logger()) }),
		fx.Provide(
			func() context.Context { return ctx },
			func() Config { return cfg },

			func(folder Folder, fs afero.Fs, tools tools.Tools) (string, error) {
				folderPath, err := filepath.Abs(string(folder))
				if err != nil {
					return "", fmt.Errorf("invalid path: %q: %w", folderPath, err)
				}

				err = fs.MkdirAll(string(folder), 0o755)
				if err != nil && !errors.Is(err, os.ErrExist) {
					return "", fmt.Errorf("failed to create the %s: %w", folderPath, err)
				}

				if fs.Name() == afero.NewMemMapFs().Name() {
					tools.Logger().Info("Load data from memory")
				} else {
					tools.Logger().Info(fmt.Sprintf("Load data from %s", folder))
				}

				return folderPath, nil
			},

			// Tools
			fx.Annotate(tools.NewToolbox, fx.As(new(tools.Tools))),
			fx.Annotate(html.NewRenderer, fx.As(new(html.Writer))),
			sqlstorage.Init,
			auth.NewAuthenticator,

			// Services
			fx.Annotate(users.Init, fx.As(new(users.Service))),
			fx.Annotate(websessions.Init, fx.As(new(websessions.Service))),
			fx.Annotate(sysstats.Init),

			// Middlewares
			middlewares.NewBootstrapMiddleware,

			// HTTP handlers
			AsRoute(assets.NewHTTPHandler),
			AsRoute(utilities.NewHTTPHandler),

			// Web Pages
			AsRoute(auth.NewLoginPage),
			AsRoute(auth.NewBootstrapPage),
			AsRoute(sysstatstmpl.NewSysstatsPage),

			// HTTP Router / HTTP Server
			router.InitMiddlewares,
			fx.Annotate(router.NewServer, fx.ParamTags(`group:"routes"`)),
		),

		fx.Invoke(migrations.Run),

		// Start the tasks-runner
		fx.Invoke(func(svc *sysstats.SystatsCron, lc fx.Lifecycle, tools tools.Tools) {
			cronSvc := cron.New(svc.Name(), svc.Duration(), tools, svc)
			cronSvc.FXRegister(lc)
		}),

		invoke,
	)

	return app
}
