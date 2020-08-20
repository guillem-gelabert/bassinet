package bassinet

import (
	"net/http"
)

const (
	deny = iota
	sameOrigin
)

// Frameguard sets the X-Frame-Options header to mitigate clickjacking attacks
func Frameguard(o int) (Middleware, error) {
	options := map[int]string{
		0: "DENY",
		1: "SAMEORIGIN",
	}

	headerValue, ok := options[o]

	if !ok {
		headerValue = options[1]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Frame-Options", headerValue)

			next.ServeHTTP(w, r)
		})
	}, nil
}
