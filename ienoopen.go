package bassinet

import "net/http"

// IeNoOpen sets X-Download-Options to "noopen" to prevent IE users
// to execute downloads in your site's context
func IeNoOpen() (Middleware, error) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Download-Options", "noopen")

			next.ServeHTTP(w, r)
		})
	}, nil
}
