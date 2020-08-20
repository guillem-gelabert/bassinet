package bassinet

import "net/http"

// DNSPrefetchControl sets X-DNS-Prefetch-Control to help prevent
// an eavesdropperfrom inferring the host names of hyperlinks
// that appear in HTTPS pages based on DNS prefetch traffic
func DNSPrefetchControl(allowed bool) Middleware {
	headerValue := "off"
	if allowed {
		headerValue = "on"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-DNS-Prefetch-Control", headerValue)

			next.ServeHTTP(w, r)
		})
	}
}
