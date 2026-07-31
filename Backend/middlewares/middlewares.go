package middlewares

import "net/http"

// MaxBodyBytes limits the maximum size of request bodies across handlers.
func MaxBodyBytes(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body = http.MaxBytesReader(res, req.Body, maxBytes)
			}
			
			next.ServeHTTP(res, req)
		})
	}
}