package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/tools/sqlstorage"
	"github.com/Peltoche/zapette/internal/web/html"
)

type BootstrapMiddleware struct {
	users         users.Service
	html          html.Writer
	isBootstraped bool
}

func NewBootstrapMiddleware(users users.Service, html html.Writer) *BootstrapMiddleware {
	return &BootstrapMiddleware{
		users:         users,
		html:          html,
		isBootstraped: false,
	}
}

func (m *BootstrapMiddleware) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isBootstraped {
			res, err := m.users.GetAll(r.Context(), &sqlstorage.PaginateCmd{Limit: 1})
			if err != nil {
				m.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to fetch the users: %w", err))
				return
			}

			m.isBootstraped = len(res) != 0
		}

		switch {
		case !m.isBootstraped && r.URL.Path != "/web/bootstrap": // Not bootstraped and not in the bootstrap process: redirect
			http.Redirect(w, r, "/web/bootstrap", http.StatusSeeOther)
		case m.isBootstraped && r.URL.Path == "/web/bootstrap": // Already bootstraped and in the bootstrap process: redirect
			http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		default:
			next.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
