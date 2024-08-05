package sysstats

import (
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
}

func (h *SysstatsPage) printListPage(w http.ResponseWriter, r *http.Request) {
	user, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	stats, err := h.sysstats.GetLast5mn(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the latest 5mn stats: %w", err))
		return
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

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &sysstatstmpl.SysstatsPageTmpl{
		Latest:      stats[len(stats)-1],
		User:        user,
		Labels:      labels,
		MemoryUsed:  memoryUsed,
		MemoryTotal: memoryTotal,
		SwapUsed:    swapUsed,
		CacheBuffer: cacheBuffer,
	})
}
