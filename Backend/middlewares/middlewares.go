package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"ScheduleFlow/Backend/services"
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

// CheckNotVerified checks if the user is not verified. If the user is verified, it redirects to the home page. 
// Prevents users from accessing the verification page if they are already verified.
func CheckNotVerified(sm *scs.SessionManager, emailVerificationService services.EmailVerificationService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")

            //query the database to check if the user is verified
            result := emailVerificationService.GetEmailVerificationEntry(userID)

            if result.Err != nil {
                http.Error(res, result.Err.Error(), http.StatusInternalServerError)
                return
            }

            // If the user has already verified their email (no email verification found), redirect to the home page
            if result.StatusCode == http.StatusNotFound {

                if req.Method == http.MethodPost {
                    res.Header().Set("HX-Redirect", "/home")
                    res.WriteHeader(http.StatusOK)
                }else{
                    http.Redirect(res, req, "/home", http.StatusSeeOther)
                }

                return
            }

            next.ServeHTTP(res, req)
        })
    }
}

// CheckVerified checks if the user is verified. If not, it redirects to the verification page. 
// Prevents users from accessing certain routes if they haven't verified their email yet.
func CheckVerified(sm *scs.SessionManager, emailVerificationService services.EmailVerificationService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")

            //query the database to check if the user is verified
            result := emailVerificationService.GetEmailVerificationEntry(userID)

            if result.Err != nil {
                http.Error(res, result.Err.Error(), result.StatusCode)
                return
            }

            // If the user has not verified their email (email verification found), redirect to the verification page
            if result.StatusCode == http.StatusOK {
                
                if req.Method == http.MethodPost {
                    res.Header().Set("HX-Redirect", "/verification")
                    res.WriteHeader(http.StatusOK)
                }else{
                    http.Redirect(res, req, "/verification", http.StatusSeeOther)
                }
                
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

                if req.Method == http.MethodPost {
                    res.Header().Set("HX-Redirect", "/home")
                    res.WriteHeader(http.StatusOK)
                }else{
                    http.Redirect(res, req, "/home", http.StatusSeeOther)
                }

                return
            }

            next.ServeHTTP(res, req)
        })
    }
}
