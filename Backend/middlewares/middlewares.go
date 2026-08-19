package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            userID := sm.GetInt(r.Context(), "userID")

            if userID == 0 {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}