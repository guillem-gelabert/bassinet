package bassinet

import "net/http"

// HidePoweredBy removes the X-POWERED-BY header to prevent attackers
// from knowing which technology the server uses
func HidePoweredBy() (Middleware, error) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Del("X-Powered-By")
			next.ServeHTTP(w, r)
		})
	}, nil
}
