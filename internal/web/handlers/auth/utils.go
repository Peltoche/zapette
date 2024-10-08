package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/service/websessions"
	"github.com/Peltoche/zapette/internal/web/html"
)

type AccessType int

const (
	AdminOnly AccessType = iota
	AnyUser
)

type Authenticator struct {
	webSessions websessions.Service
	users       users.Service
	html        html.Writer
}

func NewAuthenticator(webSessions websessions.Service, users users.Service, html html.Writer) *Authenticator {
	return &Authenticator{webSessions, users, html}
}

func (a *Authenticator) GetUserAndSession(w http.ResponseWriter, r *http.Request, access AccessType) (*users.User, *websessions.Session, bool) {
	currentSession, err := a.webSessions.GetFromReq(r)
	switch {
	case err == nil:
		break
	case errors.Is(err, websessions.ErrSessionNotFound):
		a.webSessions.Logout(r, w)
		return nil, nil, true
	case errors.Is(err, websessions.ErrMissingSessionToken):
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return nil, nil, true
	default:
		a.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to websessions.GetFromReq: %w", err))
		return nil, nil, true
	}

	user, err := a.users.GetByID(r.Context(), currentSession.UserID())
	if err != nil {
		a.html.WriteHTMLErrorPage(w, r, err)
		return nil, nil, true
	}

	if user == nil {
		_ = a.webSessions.Logout(r, w)
		return nil, nil, true
	}

	if access == AdminOnly && !user.IsAdmin() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div class="alert alert-danger role="alert">Action reserved to admins</div>`))
		return nil, nil, true
	}

	return user, currentSession, false
}

func (a *Authenticator) Logout(w http.ResponseWriter, r *http.Request) {
	a.webSessions.Logout(r, w)
}
