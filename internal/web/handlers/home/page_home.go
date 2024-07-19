package home

import (
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/systats"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/home"
	"github.com/go-chi/chi/v5"
)

type HomePage struct {
	auth    *auth.Authenticator
	html    html.Writer
	systats systats.Service
}

func NewHomePage(
	html html.Writer,
	auth *auth.Authenticator,
	systats systats.Service,
	tools tools.Tools,
) *HomePage {
	return &HomePage{
		html:    html,
		systats: systats,
		auth:    auth,
	}
}

func (h *HomePage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/", http.RedirectHandler("/web/home", http.StatusFound).ServeHTTP)
	r.Get("/web/home", h.printPage)
}

func (h *HomePage) printPage(w http.ResponseWriter, r *http.Request) {
	user, _, abort := h.auth.GetUserAndSession(w, r, auth.AnyUser)
	if abort {
		return
	}

	stats, err := h.systats.GetLatest(r.Context())
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to fetch the stats: %w", err))
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &home.HomePageTmpl{
		User: user,
		MemoryBar: home.ValueBar{
			Value: stats.Memory().UsedMemory(),
			Total: stats.Memory().TotalMemory(),
			Label: "Memory",
		},
		SwapBar: home.ValueBar{
			Value: stats.Memory().UsedSwap(),
			Total: stats.Memory().TotalSwap(),
			Label: "Swap",
		},
	})
}
