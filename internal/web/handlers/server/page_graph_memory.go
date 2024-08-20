package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/tools/ptr"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/server"
	"github.com/go-chi/chi/v5"
)

type MemoryGraphPage struct {
	html     html.Writer
	auth     *auth.Authenticator
	sysstats sysstats.Service
}

func NewMemoryGraphPage(
	html html.Writer,
	auth *auth.Authenticator,
	sysstats sysstats.Service,
) *MemoryGraphPage {
	return &MemoryGraphPage{
		html:     html,
		sysstats: sysstats,
		auth:     auth,
	}
}

func (h *MemoryGraphPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/server/memory/details", h.printMemoryGraphPage)
}

func (h *MemoryGraphPage) printMemoryGraphPage(w http.ResponseWriter, r *http.Request) {
	_, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	stats, err := h.sysstats.GetLast5mn(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the latest 5mn stats: %w", err))
		return
	}

	graphData := StatsToMemoryGraphData(stats)

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &server.SysstatsPageTmpl{
		GraphData: graphData,
	})
}

func StatsToMemoryGraphData(stats []sysstats.Stats) *server.Graph {
	memoryTotal := make([]*float64, len(stats))
	memoryUsed := make([]*float64, len(stats))
	swapUsed := make([]*float64, len(stats))
	cacheBuffer := make([]*float64, len(stats))
	labels := make([]*string, len(stats))

	for i, stat := range stats {
		if stat.IsEmpty() {
			memoryUsed[i] = nil
			continue
		}

		labels[i] = ptr.To(stat.Time().Format(time.TimeOnly))
		memoryUsed[i] = ptr.To(stat.Memory().UsedMemory().GBytes())
		memoryTotal[i] = ptr.To(stat.Memory().TotalMemory().GBytes())
		swapUsed[i] = ptr.To(stat.Memory().UsedSwap().GBytes())
		cacheBuffer[i] = ptr.To(stat.Memory().BufCache().GBytes())
	}

	return &server.Graph{
		Type: "line",
		Data: server.Data{
			Labels: labels,
			Datasets: []server.Dataset{
				{
					Label:       "Total",
					Data:        memoryTotal,
					ShowLine:    true,
					BorderColor: "red",
					SteppedLine: true,
					BorderWidth: 1,
					PointRadius: 0,
				},
				{
					Label:       "RAM",
					Data:        memoryUsed,
					ShowLine:    true,
					BorderColor: "black",
					BorderWidth: 1,
					PointRadius: 0,
				},
				{
					Label:       "Cache + Buffer",
					Data:        cacheBuffer,
					ShowLine:    true,
					BorderColor: "blue",
					BorderWidth: 1,
					PointRadius: 0,
				},
				{
					Label:       "Swap",
					Data:        swapUsed,
					ShowLine:    true,
					BorderColor: "#bf1b00",
					BorderWidth: 2,
					PointRadius: 0,
				},
			},
		},
	}
}
