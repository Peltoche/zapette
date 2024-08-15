package sysstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/ptr"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	sysstatstmpl "github.com/Peltoche/zapette/internal/web/html/templates/sysstats"
	"github.com/go-chi/chi/v5"
)

type SysstatsPage struct {
	auth     *auth.Authenticator
	html     html.Writer
	sysstats sysstats.Service
}

func NewSysstatsPage(
	html html.Writer,
	auth *auth.Authenticator,
	sysstats sysstats.Service,
	tools tools.Tools,
) *SysstatsPage {
	return &SysstatsPage{
		html:     html,
		sysstats: sysstats,
		auth:     auth,
	}
}

func (h *SysstatsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/", http.RedirectHandler("/web/sysstats", http.StatusFound).ServeHTTP)
	r.Get("/web/sysstats", h.printListPage)
	r.Get("/web/sysstats/data.json", h.servGraphData)
	r.Get("/web/sysstats/sse", h.sse)
}

func (h *SysstatsPage) servGraphData(w http.ResponseWriter, r *http.Request) {
	_, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	graphData, err := h.getGraphData(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(graphData)
}

func (h *SysstatsPage) printListPage(w http.ResponseWriter, r *http.Request) {
	_, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	graphData, err := h.getGraphData(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &sysstatstmpl.SysstatsPageTmpl{
		GraphData: graphData,
	})
}

func (h *SysstatsPage) sse(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	eventCh := h.sysstats.Watch(r.Context())

	// Send data to the client
	for range eventCh {
		graphData, err := h.getGraphData(r.Context())
		if err != nil {
			h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the graph data: %w", err))
			return
		}

		rawData, err := json.Marshal(graphData)
		if err != nil {
			h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to encode the graph data: %w", err))
			return
		}

		fmt.Fprintf(w, "event: RefreshGraph\ndata: %s\n\n", rawData)
		w.(http.Flusher).Flush()
	}
}

func (h *SysstatsPage) getGraphData(ctx context.Context) (*sysstatstmpl.Graph, error) {
	stats, err := h.sysstats.GetLast5mn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest 5mn stats: %w", err)
	}

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

	return &sysstatstmpl.Graph{
		Type: "line",
		Data: sysstatstmpl.Data{
			Labels: labels,
			Datasets: []sysstatstmpl.Dataset{
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
	}, nil
}
