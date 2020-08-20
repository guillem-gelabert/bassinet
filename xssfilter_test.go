package bassinet

import (
	"testing"
)

func TestXssFilter(t *testing.T) {
	xssFilter, err := XSSFilter()
	if err != nil {
		t.Fatal(err)
	}
	check(t, xssFilter, "X-XSS-PROTECTION", "0")
}
