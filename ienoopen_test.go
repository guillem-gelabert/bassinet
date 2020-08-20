package bassinet

import "testing"

func TestIeNoOpen(t *testing.T) {
	ieNoOpen, _ := IeNoOpen()
	check(t, ieNoOpen, "X-Download-Options", "noopen")
}
