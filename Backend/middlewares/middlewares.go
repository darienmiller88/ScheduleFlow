package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"ScheduleFlow/Backend/services"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

//PreventResendCode prevents users from spamming the resend verification code button. It checks if the user has
//already requested a code within the last 15 minutes. If they have, it returns a 429 Too Many Requests
//status code and an error message. If not, it allows the request to proceed and sets a rate limit for the user.
func PreventResendCode(sm *scs.SessionManager, emailVerificationService services.EmailVerificationService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            userID := sm.GetInt(req.Context(), "userID")

            //query the database to check if the user is verified
            result := emailVerificationService.GetEmailVerificationEntry(userID)

            if result.Err != nil {
                http.Error(res, result.Err.Error(), http.StatusInternalServerError)
                return
            }

            // If the user has already requested a code within the last 15 minutes, return a 429 Too Many Requests status code
            if timeUntilExpires := time.Until(result.ResultData.ExpiresAt); timeUntilExpires > 0 {
                http.Error(res, fmt.Sprintf("You have already requested a verification code. Please wait %d minutes before requesting again.", int(timeUntilExpires.Minutes())), http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(res, req)
        })
    }
}


// ClientIPKey is the rate-limit key. middleware.GetClientIP reads the IP
// resolved in step 1; httprate.CanonicalizeIP buckets IPv6 clients by /64.
func ClientIPKey(r *http.Request) (string, error) {
	return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
}

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
                fmt.Println("User is authenticated, redirecting to home page:", userID)

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
