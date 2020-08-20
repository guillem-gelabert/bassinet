package bassinet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Middleware function
type Middleware = func(next http.Handler) http.Handler

func check(t *testing.T, m Middleware, headerName, expectedValue string) *http.Response {
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

	actualValue := rs.Header.Get(headerName)
	if actualValue != expectedValue {
		t.Errorf("expected \n%q\ngot \n%q", expectedValue, actualValue)
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

func pretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}
