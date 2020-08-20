package bassinet

import (
	"testing"
)

func TestXPoweredBy(t *testing.T) {
	hidePoweredBy, err := HidePoweredBy()
	if err != nil {
		t.Fatal(err)
	}
	check(t, hidePoweredBy, "X-POWERED-BY", "")
}
