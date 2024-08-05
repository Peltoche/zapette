package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/service/websessions"
	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/Peltoche/zapette/internal/tools/router"
	"github.com/Peltoche/zapette/internal/tools/secret"
	"github.com/Peltoche/zapette/internal/tools/uuid"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/Peltoche/zapette/internal/web/html/templates/auth"
	"github.com/go-chi/chi/v5"
)

const cookieLifeTime = time.Hour * 24 * 365

type LoginPage struct {
	webSessions websessions.Service
	uuid        uuid.Service
	html        html.Writer
	users       users.Service
	clock       clock.Clock
}

func NewLoginPage(
	html html.Writer,
	webSessions websessions.Service,
	users users.Service,
	tools tools.Tools,
) *LoginPage {
	return &LoginPage{
		html:        html,
		webSessions: webSessions,
		users:       users,
		uuid:        tools.UUID(),
		clock:       tools.Clock(),
	}
}

func (h *LoginPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/login", h.printPage)
	r.Post("/web/login", h.applyLogin)
}

func (h *LoginPage) printPage(w http.ResponseWriter, r *http.Request) {
	currentSession, _ := h.webSessions.GetFromReq(r)

	if currentSession != nil {
		h.chooseRedirection(w, r)
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &auth.LoginPageTmpl{})
}

func (h *LoginPage) applyLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := auth.LoginPageTmpl{}

	tmpl.UsernameContent = r.FormValue("username")

	user, err := h.users.Authenticate(r.Context(), r.FormValue("username"), secret.NewText(r.FormValue("password")))
	var status int
	switch {
	case err == nil:
		// continue
	case errors.Is(err, users.ErrInvalidUsername):
		tmpl.UsernameError = "User doesn't exists"
		status = http.StatusBadRequest
	case errors.Is(err, users.ErrInvalidPassword):
		tmpl.PasswordError = "Invalid password"
		status = http.StatusBadRequest
	default:
		h.html.WriteHTMLErrorPage(w, r, err)
		return
	}

	if err != nil {
		h.html.WriteHTMLTemplate(w, r, status, &tmpl)
		return
	}

	session, err := h.webSessions.Create(r.Context(), &websessions.CreateCmd{
		UserID:     user.ID(),
		UserAgent:  r.Header.Get("User-Agent"),
		RemoteAddr: r.RemoteAddr,
	})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to create the websession: %w", err))
		return
	}

	var expirationDate time.Time
	if r.FormValue("remember") != "" {
		expirationDate = h.clock.Now().Add(cookieLifeTime)
	}

	c := http.Cookie{
		Name:     "session_token",
		Value:    session.Token().Raw(),
		Expires:  expirationDate,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &c)

	h.chooseRedirection(w, r)
}

func (h *LoginPage) chooseRedirection(w http.ResponseWriter, r *http.Request) {
	// NOTE: The oauth dance is not implemented yet
	http.Redirect(w, r, "/web/sysstats", http.StatusFound)

	// TODO: Uncomment this once the oauth2 is implemented
	//
	//  	var client *oauthclients.Client
	//  	clientID, err := h.uuid.Parse(r.FormValue("client_id"))
	//  	if err == nil {
	//  		client, err = h.clients.GetByID(r.Context(), clientID)
	//  		if err != nil {
	//  			reqID, ok := r.Context().Value(middleware.RequestIDKey).(string)
	//  			if !ok {
	//  				reqID = "????"
	//  			}
	//
	//  			h.html.WriteHTMLTemplate(w, r, http.StatusBadRequest, &auth.ErrorPageTmpl{
	//  				ErrorMsg:  "Oauth client not found",
	//  				RequestID: reqID,
	//  			})
	//  			return
	//  		}
	//  	}
	//
	//  	switch {
	//  	case client == nil:
	//  		http.Redirect(w, r, "/", http.StatusFound)
	//  	case client.SkipValidation():
	//  		http.Redirect(w, r, "/auth/authorize", http.StatusFound)
	//  	default:
	//  		http.Redirect(w, r, "/consent?"+r.Form.Encode(), http.StatusFound)
	//  	}
}
