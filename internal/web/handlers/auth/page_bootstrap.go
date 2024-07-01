package auth

import (
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/tools/secret"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/auth"
	"github.com/go-chi/chi/v5"
)

type BootstrapPage struct {
	html  html.Writer
	users users.Service
}

func NewBootstrapPage(html html.Writer, users users.Service) *BootstrapPage {
	return &BootstrapPage{
		html:  html,
		users: users,
	}
}

func (h *BootstrapPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/bootstrap", h.printPage)
	r.Post("/web/bootstrap", h.postForm)
}

func (h *BootstrapPage) printPage(w http.ResponseWriter, r *http.Request) {
	// if h.masterkey.IsMasterKeyLoaded() {
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &auth.BootstrapPageTmpl{})
}

func (h *BootstrapPage) postForm(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := secret.NewText(r.FormValue("password"))
	confirm := secret.NewText(r.FormValue("confirm"))

	if confirm != password {
		h.html.WriteHTMLTemplate(w, r, http.StatusOK, &auth.BootstrapPageTmpl{
			Username:      username,
			Password:      password,
			PasswordError: "",
			ConfirmError:  "not identical",
		})
		return
	}

	if len(password.Raw()) < 8 {
		h.html.WriteHTMLTemplate(w, r, http.StatusOK, &auth.BootstrapPageTmpl{
			Username:      username,
			Password:      password,
			PasswordError: "too short",
			ConfirmError:  "",
		})
		return
	}

	_, err := h.users.Bootstrap(r.Context(), &users.BootstrapCmd{
		Username: username,
		Password: password,
	})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to boostrap the users: %w", err))
		return
	}

	http.Redirect(w, r, "/web/login", http.StatusFound)
}
