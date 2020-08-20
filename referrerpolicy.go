package bassinet

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// PolicyNoReferrer no-referrer policy
	PolicyNoReferrer = iota
	// PolicyNoReferrerWhenDowngrade no-referrer-when-downgrade policy
	PolicyNoReferrerWhenDowngrade
	// PolicySameOrigin same-origin policy
	PolicySameOrigin
	// PolicyOrigin origin policy
	PolicyOrigin
	// PolicyStrictOrigin strict-origin policy
	PolicyStrictOrigin
	// PolicyOriginWhenCrossOrigin origin-when-cross-origin policy
	PolicyOriginWhenCrossOrigin
	// PolicyStrictOriginWhenCrossOrigin strict-origin-when-cross-origin policy
	PolicyStrictOriginWhenCrossOrigin
	// PolicyUnsafeURL unsafe-url policy
	PolicyUnsafeURL
)

// ReferrerPolicy sets the Referrer-Policy HTTP header to let authors control how browsers set the Referer header.
func ReferrerPolicy(policies []int) (Middleware, error) {
	permittedPolicies := map[int]string{
		0: "no-referrer",
		1: "no-referrer-when-downgrade",
		2: "same-origin",
		3: "origin",
		4: "strict-origin",
		5: "origin-when-cross-origin",
		6: "strict-origin-when-cross-origin",
		7: "unsafe-url",
	}

	policies = unique(append(policies, PolicyNoReferrer))
	headerValue := []string{}

	for _, p := range policies {
		policy, ok := permittedPolicies[p]
		headerValue = append(headerValue, policy)

		if !ok {
			return nil, fmt.Errorf("Referrer-Policy does not support %d", p)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Referrer-Policy", strings.Join(headerValue, ","))

			next.ServeHTTP(w, r)
		})
	}, nil
}

func unique(ns []int) (r []int) {
	seen := make(map[int]bool)
	for _, n := range ns {
		if _, ok := seen[n]; !ok {
			r = append(r, n)
		}
		seen[n] = true
	}
	return r
}
