package sysstats

import (
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/systats"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/sysstats"
	sysstatstmpl "github.com/Peltoche/zapette/internal/web/html/templates/sysstats"
	"github.com/go-chi/chi/v5"
)

type SysstatsPage struct {
	auth    *auth.Authenticator
	html    html.Writer
	systats systats.Service
}

func NewSysstatsPage(
	html html.Writer,
	auth *auth.Authenticator,
	systats systats.Service,
	tools tools.Tools,
) *SysstatsPage {
	return &SysstatsPage{
		html:    html,
		systats: systats,
		auth:    auth,
	}
}

func (h *SysstatsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/", http.RedirectHandler("/web/sysstats", http.StatusFound).ServeHTTP)
	r.Get("/web/sysstats", h.printPage)
}

func (h *SysstatsPage) printPage(w http.ResponseWriter, r *http.Request) {
	user, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	stats, err := h.systats.GetLatest(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to fetch the stats: %w", err))
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &sysstatstmpl.SysstatsPageTmpl{
		User: user,
		MemoryBar: sysstats.ValueBar{
			Value: stats.Memory().UsedMemory(),
			Total: stats.Memory().TotalMemory(),
			Label: "Memory",
		},
		SwapBar: sysstats.ValueBar{
			Value: stats.Memory().UsedSwap(),
			Total: stats.Memory().TotalSwap(),
			Label: "Swap",
		},
	})
}
