package bassinet

import "testing"

func TestDNSPrefetchControl(t *testing.T) {
	testCases := []struct {
		desc     string
		target   bool
		expected string
	}{
		{
			desc:     "disallowed",
			target:   false,
			expected: "off",
		},
		{
			desc:     "allowed",
			target:   true,
			expected: "on",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			dnsPrefetchControl := DNSPrefetchControl(tC.target)

			check(t, dnsPrefetchControl, "X-DNS-Prefetch-Control", tC.expected)
		})
	}
}
