package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"

    "ScheduleFlow/Backend/repositories"
)

// RequireAuth checks if the user is authenticated. If not, it redirects to the login page.
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

// RequireVerification checks if the user is verified. If so, it redirects to the home page.
func RequireVerification(sm *scs.SessionManager, db *sqlx.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")
            repository := repositories.NewEmailVerificationRepository(db)

            //query the database to check if the user is verified
            result := repository.GetEmailVerification(userID)

            if result.Err != nil {
                http.Error(res, result.Err.Error(), http.StatusInternalServerError)
                return
            }

            // If the user has already verified their email (no email verification found), redirect to the home page
            if result.StatusCode == http.StatusNotFound {
                http.Redirect(res, req, "/home", http.StatusSeeOther)
                return
            }

            next.ServeHTTP(res, req)
        })
    }
}

// SendBackToHome checks if the user is authenticated. If so, it redirects to the home page.
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
