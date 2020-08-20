package bassinet

import (
	"fmt"
	"net/http"
	"strings"
)

// ExpectCtOptions options for the Expect-CT header
type ExpectCtOptions struct {
	enforce   bool
	maxAge    int32
	reportURI string
}

// ExpectCt sets the Expect-CT header to expect Certificate Transparency
func ExpectCt(o ExpectCtOptions) (Middleware, error) {
	if o.maxAge < 0 {
		return nil, fmt.Errorf("max-age must be a positive integer")
	}

	directives := []string{fmt.Sprintf("max-age=%d", o.maxAge)}
	if o.enforce {
		directives = append(directives, "enforce")
	}
	if o.reportURI != "" {
		directives = append(directives, `report-uri=`+`"`+o.reportURI+`"`)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Expect-CT", strings.Join(directives, ", "))

			next.ServeHTTP(w, r)
		})
	}, nil
}
