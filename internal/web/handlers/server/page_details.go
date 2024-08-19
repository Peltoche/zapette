package server

import (
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/sysstats"
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
}

func NewDetailsPage(
	html html.Writer,
	auth *auth.Authenticator,
	sysstats sysstats.Service,
) *DetailsPage {
	return &DetailsPage{
		html:     html,
		sysstats: sysstats,
		auth:     auth,
	}
}

func (h *DetailsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/server", h.printServerPage)
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
		Stats: latest,
	})
}
