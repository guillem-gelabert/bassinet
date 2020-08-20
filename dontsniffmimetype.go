package bassinet

import "net/http"

// DontSniffMimetype sets X-Content-Type-Options to nosniff to avoid browsers from
// executing textfiles
func DontSniffMimetype() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")

			next.ServeHTTP(w, r)
		})
	}
}
