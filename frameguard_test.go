package bassinet

import "testing"

func TestFrameguard(t *testing.T) {
	testCases := []struct {
		desc     string
		target   int
		expected string
	}{
		{
			desc:     "with DENY",
			target:   0,
			expected: "DENY",
		},
		{
			desc:     "with SAMEORIGIN",
			target:   1,
			expected: "SAMEORIGIN",
		},
		{
			desc:     "defaults to SAMEORIGIN",
			target:   10,
			expected: "SAMEORIGIN",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			frameguard, err := Frameguard(tC.target)
			if err != nil {
				t.Fatal(err)
			}
			check(t, frameguard, "X-Frame-Options", tC.expected)
		})
	}
}
