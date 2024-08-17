package router

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/misc"
	"github.com/coreos/go-systemd/daemon"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/afero"
	"go.uber.org/fx"
)

type API struct{}

type Config struct {
	Addr      string
	CertFile  string
	KeyFile   string
	HostNames []string
	TLS       bool
	Secure    bool
}

type Registerer interface {
	Register(r chi.Router, mids *Middlewares)
}

type ConnectionCloser interface {
	CloseOpenConnections()
}

func NewServer(
	routes []Registerer,
	cfg Config,
	lc fx.Lifecycle,
	mids *Middlewares,
	tools tools.Tools,
	fs afero.Fs,
	writer html.Writer,
) (*API, *http.Server, error) {
	handler, err := createHandler(cfg, routes, mids, writer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the listener: %w", err)
	}

	httpLogger := slog.NewLogLogger(tools.Logger().Handler(), slog.LevelError)

	srv := &http.Server{
		Addr:     cfg.Addr,
		Handler:  handler,
		ErrorLog: httpLogger,
	}

	srv.RegisterOnShutdown(func() {
		for _, route := range routes {
			if r, ok := route.(ConnectionCloser); ok {
				r.CloseOpenConnections()
			}
		}
	})

	if cfg.TLS {
		cert, err := afero.ReadFile(fs, cfg.CertFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load the TLS certification file: %w", err)
		}

		key, err := afero.ReadFile(fs, cfg.KeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load the TLS key file: %w", err)
		}

		certif, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate the X509 key pair: %w", err)
		}

		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certif},
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			ln, err := net.Listen("tcp", cfg.Addr)
			if err != nil {
				return err
			}

			tools.Logger().Info("start listening", slog.String("host", ln.Addr().String()), slog.Int("routes", len(handler.Routes())))
			if cfg.TLS {
				go srv.ServeTLS(ln, "", "")
			} else {
				go srv.Serve(ln)
			}

			for _, route := range handler.Routes() {
				tools.Logger().Debug("expose endpoint", slog.String("host", ln.Addr().String()), slog.String("route", route.Pattern))
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	daemon.SdNotify(false, daemon.SdNotifyReady)

	return &API{}, srv, nil
}

func createHandler(cfg Config, routes []Registerer, mids *Middlewares, writer html.Writer) (chi.Router, error) {
	r := chi.NewMux()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writer.WriteHTMLTemplate(w, r, http.StatusNotFound, &misc.NotFoundPageTmpl{})
		// res.WriteJSON(w http.ResponseWriter, r *http.Request, statusCode int, res any)
		// http.Redirect(w, r, "", http.StatusFound)
	})

	if cfg.Secure {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Strict-Transport-Security", "max-age=15768000; preload")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("Referrer-Policy", "same-origin")
				w.Header().Set("Permissions-Policy", "")

				next.ServeHTTP(w, r)
			})
		})
	}

	r.Use(mids.CORS)
	r.Use(middleware.RequestID)

	for _, svc := range routes {
		svc.Register(r, mids)
	}

	return r, nil
}
