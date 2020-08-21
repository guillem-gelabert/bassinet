package bassinet

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestCSP(t *testing.T) {
	testCases := []struct {
		desc           string
		options        CSPOptions
		expectedHeader string
		expectedValues map[string][]string
	}{
		{
			desc: "it sets defaults and report only",
			options: CSPOptions{
				reportOnly: true,
			},
			expectedHeader: "Content-Security-Policy-Report-Only",
			expectedValues: defaultDirectives,
		},
		{
			desc: "it sets policies",
			options: CSPOptions{
				directives: map[string][]string{
					"default-src": {"'none'"},
					"script-src":  {"https://cdn.mybank.net"},
					"style-src":   {"https://cdn.mybank.net"},
					"img-src":     {"https://cdn.mybank.net"},
					"connect-src": {"https://api.mybank.com"},
					"child-src":   {"'self'"},
				},
			},
			expectedHeader: "Content-Security-Policy",
			expectedValues: parseDirectives(`default-src 'none'; script-src https://cdn.mybank.net; style-src https://cdn.mybank.net; img-src https://cdn.mybank.net; connect-src https://api.mybank.com; child-src 'self'`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			csp, err := CSP(tC.options)
			if err != nil {
				t.Fatal(err)
			}

			rs := checkCSP(t, csp, tC.expectedHeader, tC.expectedValues)
			header := rs.Header.Get(tC.expectedHeader)
			if header == "" {
				t.Errorf("Excpected %q header to be present", tC.expectedHeader)
			}
		})
	}

	t.Run("it throws if default-src is missing", func(t *testing.T) {
		_, err := CSP(CSPOptions{
			directives: map[string][]string{
				"script-src": {"https://cdn.mybank.net"},
			},
		})
		if err == nil {
			t.Error(err)
		}
	})
}

func TestToKebab(t *testing.T) {
	testCases := []struct {
		desc     string
		target   string
		expected string
	}{
		{
			desc:     "valid camel cased string",
			target:   "abcDefGhi",
			expected: "abc-def-ghi",
		},
		{
			desc:     "valid camel cased string with numbers",
			target:   "abcDefGhi3",
			expected: "abc-def-ghi3",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := toKebab(tC.target)
			if err != nil {
				t.Fatal(err)
			}
			if actual != tC.expected {
				t.Errorf("Expected %q; got %q", tC.expected, actual)
			}
		})
	}
	t.Run("fails with non english characters", func(t *testing.T) {
		_, err := toKebab("àsd")
		if err == nil {
			t.Error("Expected to fail")
		}
	})
}

func TestNormalizeDirectives(t *testing.T) {
	testCases := []struct {
		desc         string
		target       map[string][]string
		expectedKeys []string
	}{
		{
			desc:   "with default directives",
			target: defaultDirectives,
			expectedKeys: []string{
				"default-src",
				"base-uri",
				"block-all-mixed-content",
				"font-src",
				"frame-ancestors",
				"img-src",
				"object-src",
				"script-src",
				"script-src-attr",
				"style-src",
				"upgrade-insecure-requests",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			normedDirectives, err := normalizeDirectives(tC.target)
			if err != nil {
				t.Fatal(err)
			}

			var actualKeys []string
			for key := range normedDirectives {
				actualKeys = append(actualKeys, key)
			}

			sort.Strings(actualKeys)
			sort.Strings(tC.expectedKeys)
			if !reflect.DeepEqual(tC.expectedKeys, actualKeys) {
				t.Errorf("Expected %v; got %v", tC.expectedKeys, actualKeys)
			}
		})
	}
}

func TestSerializeDirectives(t *testing.T) {
	testCases := []struct {
		desc     string
		target   map[string][]string
		expected string
	}{
		{
			desc: "",
			target: map[string][]string{
				"child-src": {
					"https://plusone.google.com",
					"https://facebook.com",
					"https://platform.twitter.com"},
			},
			expected: "child-src https://plusone.google.com https://facebook.com https://platform.twitter.com; ",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := serializeDirectives(tC.target)
			if actual != tC.expected {
				t.Errorf("Expected %q got %q", tC.expected, actual)
			}
		})
	}
}

func TestParseDirectives(t *testing.T) {
	testCases := []struct {
		desc     string
		target   string
		expected map[string][]string
	}{
		{
			desc:     "",
			target:   serializeDirectives(defaultDirectives),
			expected: defaultDirectives,
		},
		{
			desc:   "",
			target: `default-src 'none'; script-src https://cdn.mybank.net; style-src https://cdn.mybank.net; img-src https://cdn.mybank.net; connect-src https://api.mybank.com; child-src 'self'`,
			expected: map[string][]string{
				"default-src": {"'none'"},
				"script-src":  {"https://cdn.mybank.net"},
				"style-src":   {"https://cdn.mybank.net"},
				"img-src":     {"https://cdn.mybank.net"},
				"connect-src": {"https://api.mybank.com"},
				"child-src":   {"'self'"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := parseDirectives(tC.target)
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected:\n%s\ngot\n%s", pretty(tC.expected), pretty(actual))
			}
		})
	}
}

func checkCSP(t *testing.T, m Middleware, expectedHeader string, expectedValues map[string][]string) *http.Response {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	m(next).ServeHTTP(rr, r)

	rs := rr.Result()

	actualValues := rs.Header.Get(expectedHeader)
	if actualValues == "" {
		t.Errorf("expected %q to be present", expectedHeader)
	}

	parsedValues := parseDirectives(actualValues)

	if !reflect.DeepEqual(expectedValues, parsedValues) {
		t.Errorf("expected \n%q\ngot \n%q", expectedValues, parsedValues)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("expected body to be %q", "OK")
	}

	return rs
}

func parseDirectives(s string) map[string][]string {
	result := map[string][]string{}
	directives := strings.Split(s, "; ")
	for _, d := range directives {
		if d == "" {
			continue
		}
		parsedDirective := strings.Split(d, " ")

		if parsedDirective[1] != "" {
			result[strings.TrimSpace(parsedDirective[0])] = parsedDirective[1:]
		} else {
			result[parsedDirective[0]] = []string{}
		}
	}
	return result
}
