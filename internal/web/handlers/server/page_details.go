package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/sysinfos"
	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/server"
	"github.com/go-chi/chi/v5"
)

type DetailsPage struct {
	html     html.Writer
	auth     *auth.Authenticator
	sysstats sysstats.Service
	sysinfos sysinfos.Service
	logger   *slog.Logger
	closeCh  chan struct{}
}

func NewDetailsPage(
	html html.Writer,
	tools tools.Tools,
	auth *auth.Authenticator,
	sysinfos sysinfos.Service,
	sysstats sysstats.Service,
) *DetailsPage {
	return &DetailsPage{
		html:     html,
		sysstats: sysstats,
		sysinfos: sysinfos,
		auth:     auth,
		logger:   tools.Logger().With(slog.String("source", "server-details-sse")),
		closeCh:  make(chan struct{}, 1),
	}
}

func (h *DetailsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/server", h.printServerPage)
	r.Get("/web/server/sse", h.sse)
}

func (h *DetailsPage) printServerPage(w http.ResponseWriter, r *http.Request) {
	_, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	latest, err := h.sysstats.GetLatest(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the latest stat: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &server.DetailsPageTmpl{
		Stats:    latest,
		SysInfos: h.sysinfos.GetInfos(r.Context()),
	})
}

func (h *DetailsPage) sse(w http.ResponseWriter, r *http.Request) {
	type refreshPage struct {
		PercentageUsedMemory      int    `json:"percentageUsedMemory"`
		PercentageAvailableMemory int    `json:"percentageAvailableMemory"`
		TotalMemory               string `json:"totalMemory"`
	}

	ctx := r.Context()

	_, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	eventCh := h.sysstats.Watch(ctx)

	for {
		select {
		case <-eventCh:
			// continue
		case <-h.closeCh:
			return
		}

		latest, err := h.sysstats.GetLatest(ctx)
		if err != nil {
			h.logger.Error("failed to get the latest stat", slog.String("error", err.Error()))
			return
		}

		rawData, err := json.Marshal(&refreshPage{
			PercentageUsedMemory:      latest.Memory().PercentageUsedMemory(),
			PercentageAvailableMemory: latest.Memory().PercentageAvailableMemory(),
			TotalMemory:               latest.Memory().TotalMemory().HR(),
		})
		if err != nil {
			h.logger.Error("failed to marshal the latest stat", slog.String("error", err.Error()))
			continue
		}

		fmt.Fprintf(w, "event: LatestStat\ndata: %s\n\n", rawData)
		w.(http.Flusher).Flush()
	}
}

func (h *DetailsPage) CloseOpenConnections() {
	h.logger.Info("close open connections")
	close(h.closeCh)
}
