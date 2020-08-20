package bassinet

import (
	"fmt"
	"net/http"
)

const (
	none = iota
	masterOnly
	byContentType
	all
)

// PermittedCrossDomainPolicies sets X-Permitted-Cross-Domain-Policies header
// to tell some web clients your domain's policy for loading cross-domain content.
func PermittedCrossDomainPolicies(p int) (Middleware, error) {

	permittedPolicies := map[int]string{
		0: "none",
		1: "master-only",
		2: "by-content-type",
		3: "all",
	}

	headerValue, ok := permittedPolicies[p]
	if !ok {
		return nil, fmt.Errorf("X-Permitted-Cross-Domain-Policies does not support %d", p)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Permitted-Cross-Domain-Policies", headerValue)

			next.ServeHTTP(w, r)
		})
	}, nil
}
