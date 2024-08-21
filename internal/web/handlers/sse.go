package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/handlers/server"
	"github.com/go-chi/chi/v5"
)

type SSEPage struct {
	auth     *auth.Authenticator
	sysstats sysstats.Service
	logger   *slog.Logger
	closeCh  chan struct{}
}

func NewSSEPage(
	auth *auth.Authenticator,
	tools tools.Tools,
	sysstats sysstats.Service,
) *SSEPage {
	return &SSEPage{
		auth:     auth,
		sysstats: sysstats,
		logger:   tools.Logger().With(slog.String("source", "SSE")),
		closeCh:  make(chan struct{}, 1),
	}
}

func (p *SSEPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/sse", p.sse)
}

func (p *SSEPage) sse(w http.ResponseWriter, r *http.Request) {
	_, _, abort := p.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	p.listenSysstatsEvents(r.Context(), w)
}

func (p *SSEPage) listenSysstatsEvents(ctx context.Context, w http.ResponseWriter) {
	eventCh := p.sysstats.Watch(ctx)

	// Send data to the client
	for {
		select {
		case <-eventCh:
			// continue
		case <-p.closeCh:
			return
		}

		stats, err := p.sysstats.GetLast5mn(ctx)
		if err != nil {
			p.logger.Error("failed to get the 5mn stats", slog.String("error", err.Error()))
			return
		}

		graphData := server.StatsToMemoryGraphData(stats)

		rawData, err := json.Marshal(graphData)
		if err != nil {
			p.logger.Error("failed to marshal the graph data", slog.String("error", err.Error()))
			continue
		}

		fmt.Fprintf(w, "event: RefreshGraph\ndata: %s\n\n", rawData)
		w.(http.Flusher).Flush()
	}
}

func (p *SSEPage) CloseOpenConnections() {
	p.logger.Warn("close open connections")
	close(p.closeCh)
}
