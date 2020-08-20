package bassinet

import "testing"

func TestStrictTransportSecurity(t *testing.T) {
	testCases := []struct {
		desc     string
		target   StrictTransportOptions
		expected string
	}{
		{
			desc:     `defaults to max-age 180 days and "includeSubDomains"`,
			target:   StrictTransportOptions{},
			expected: "max-age=15552000; includeSubDomains",
		},
		{
			desc: "with two years as max-age and excluding subdomains",
			target: StrictTransportOptions{
				maxAge:            60 * 60 * 24 * 365 * 2,
				excludeSubdomains: true,
			},
			expected: "max-age=63072000",
		},
		{
			desc: "with default max-age, including subdomains and preload",
			target: StrictTransportOptions{
				preload: true,
			},
			expected: "max-age=15552000; includeSubDomains; preload",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			strictTransportSecurity, err := StrictTransportSecurity(tC.target)
			if err != nil {
				t.Fatal(err)
			}

			check(t, strictTransportSecurity, "Strict-Transport-Security", tC.expected)
		})
	}
}
