package bassinet

import "testing"

func TestDontSniffMimetype(t *testing.T) {
	dontSniffMimetype := DontSniffMimetype()
	check(t, dontSniffMimetype, "X-Content-Type-Options", "nosniff")
}
