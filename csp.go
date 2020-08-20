package bassinet

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

// CSPOptions takes directives as a map of strings and a flag to
// set the Content-Security-Policy-Report-Only header
type CSPOptions struct {
	directives map[string][]string
	reportOnly bool
}

var defaultDirectives map[string][]string = map[string][]string{
	"default-src":               {"'self'"},
	"base-uri":                  {"'self'"},
	"block-all-mixed-content":   {},
	"font-src":                  {"'self'", "https:", "data:"},
	"frame-ancestors":           {"'self'"},
	"img-src":                   {"'self'", "data:"},
	"object-src":                {"'none'"},
	"script-src":                {"'self'"},
	"script-src-attr":           {"'none'"},
	"style-src":                 {"'self'", "https:", "'unsafe-inline'"},
	"upgrade-insecure-requests": {},
}

// CSP sets the Content-Security-Policy header
func CSP(o CSPOptions) (Middleware, error) {
	directives := map[string][]string{}
	if len(o.directives) < 1 {
		directives = mergeDirectives(directives, defaultDirectives)
	}
	directives = mergeDirectives(directives, o.directives)

	normedDirectives, err := normalizeDirectives(directives)
	if err != nil {
		return nil, err
	}

	if _, ok := directives["default-src"]; !ok {
		return nil, fmt.Errorf("Content-Security-Policy needs a default-src but none was provided")
	}

	serializedDirectives := serializeDirectives(normedDirectives)

	headerName := "Content-Security-Policy"
	if o.reportOnly {
		headerName = "Content-Security-Policy-Report-Only"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(headerName, serializedDirectives)

			next.ServeHTTP(w, r)
		})
	}, nil
}

func toKebab(s string) (string, error) {
	var result string
	for _, r := range s {
		if r != rune('-') && (int(r) > 128 || !(unicode.IsLetter(r) || unicode.IsNumber(r))) {
			return "", fmt.Errorf("Input string contains invalid character %q", r)
		}
		if unicode.IsUpper(r) {
			result += "-" + string(unicode.ToLower(r))
			continue
		}
		result += string(r)
	}

	return result, nil
}

func normalizeDirectives(directives map[string][]string) (map[string][]string, error) {
	normedDirectives := map[string][]string{}

	for k, d := range directives {
		key, err := toKebab(k)
		if err != nil {
			return nil, err
		}
		normedDirectives[key] = append(normedDirectives[key], d...)
	}

	if len(normedDirectives["defaultSrc"]) == 1 {
		return nil, fmt.Errorf("defaultSrc must be set")
	}

	return normedDirectives, nil
}

func serializeDirectives(directives map[string][]string) string {
	var serialized string
	for k, v := range directives {
		serialized += fmt.Sprintf("%s %s; ", k, strings.Join(v, " "))
	}

	return serialized
}

func mergeDirectives(d1, d2 map[string][]string) map[string][]string {
	for k := range d2 {
		d1[k] = d2[k]
	}
	return d1
}
