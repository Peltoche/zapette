package home

import (
	"net/http"

	"github.com/Peltoche/zapette/internal/service/websessions"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/web/handlers/auth"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/home"
	"github.com/go-chi/chi/v5"
)

type HomePage struct {
	webSessions websessions.Service
	auth        *auth.Authenticator
	html        html.Writer
}

func NewHomePage(
	html html.Writer,
	auth *auth.Authenticator,
	tools tools.Tools,
) *HomePage {
	return &HomePage{
		html: html,
		auth: auth,
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

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &home.HomePageTmpl{
		User: user,
	})
}
