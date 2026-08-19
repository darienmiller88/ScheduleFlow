package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")

            // If userID is 0, it means the user is not authenticated, so redirect to login page
            if userID == 0 {
                http.Redirect(res, req, "/", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(res, req)
        })
    }
}

func SendBackToHome(sm *scs.SessionManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")

            // If userID is not 0, it means the user is authenticated, so redirect to home page
            if userID != 0 {
                http.Redirect(res, req, "/home", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(res, req)
        })
    }
}
