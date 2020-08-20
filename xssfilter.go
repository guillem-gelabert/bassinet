package bassinet

import "net/http"

// XSSFilter disables X-XSS-PROTECTION to avoid unintended security issues like xsleaks attack
func XSSFilter() (Middleware, error) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-XSS-Protection", "0")

			next.ServeHTTP(w, r)
		})
	}, nil
}
