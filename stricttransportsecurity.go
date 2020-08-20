package bassinet

import (
	"fmt"
	"net/http"
	"strings"
)

// StrictTransportOptions options for the Strict-Transport-Security header
type StrictTransportOptions struct {
	maxAge            int32
	excludeSubdomains bool
	preload           bool
}

// StrictTransportSecurity sets Strict-Transport-Security so that browsers
// use HTTPS for the specified period of time
func StrictTransportSecurity(o StrictTransportOptions) (Middleware, error) {
	var maxAge int32 = 15552000
	if o.maxAge > 0 {
		maxAge = o.maxAge
	}

	directives := []string{fmt.Sprintf("max-age=%d", maxAge)}

	if !o.excludeSubdomains {
		directives = append(directives, "includeSubDomains")
	}

	if o.preload {
		directives = append(directives, "preload")
	}

	headerValue := strings.Join(directives, "; ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", headerValue)

			next.ServeHTTP(w, r)
		})
	}, nil
}
