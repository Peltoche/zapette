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
	sysstatshandler "github.com/Peltoche/zapette/internal/web/handlers/sysstats"
	"github.com/go-chi/chi/v5"
)

type SSEPage struct {
	auth     *auth.Authenticator
	sysstats sysstats.Service
	logger   *slog.Logger
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
	}
}

func (p *SSEPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/sse", p.sse)
}

func (p *SSEPage) sse(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	p.listenSysstatsEvents(r.Context(), w)
}

func (p *SSEPage) listenSysstatsEvents(ctx context.Context, w http.ResponseWriter) {
	eventCh := p.sysstats.Watch(ctx)

	// Send data to the client
	for range eventCh {
		stats, err := p.sysstats.GetLast5mn(ctx)
		if err != nil {
			p.logger.Error("failed to get the 5mn stats", slog.String("error", err.Error()))
			return
		}

		graphData := sysstatshandler.StatsToGraphData(stats)

		rawData, err := json.Marshal(graphData)
		if err != nil {
			p.logger.Error("failed to marshal the graph data", slog.String("error", err.Error()))
			continue
		}

		fmt.Fprintf(w, "event: RefreshGraph\ndata: %s\n\n", rawData)
		w.(http.Flusher).Flush()
	}
}
