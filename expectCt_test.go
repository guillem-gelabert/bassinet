package bassinet

import "testing"

func TestExpectCt(t *testing.T) {
	testCases := []struct {
		desc     string
		options  ExpectCtOptions
		expected string
	}{
		{
			desc:     "defaults max-age to 0",
			options:  ExpectCtOptions{},
			expected: "max-age=0",
		},
		{
			desc: "with max age to one week",
			options: ExpectCtOptions{
				maxAge: 60 * 60 * 24 * 7,
			},
			expected: "max-age=604800",
		},
		{
			desc: "with enforcement",
			options: ExpectCtOptions{
				enforce: true,
			},
			expected: "max-age=0, enforce",
		},
		{
			desc: "with report-uri",
			options: ExpectCtOptions{
				reportURI: "http://www.domain.com",
			},
			expected: `max-age=0, report-uri="http://www.domain.com"`,
		},
		{
			desc: "with all options set",
			options: ExpectCtOptions{
				enforce:   true,
				maxAge:    60 * 60 * 24 * 7,
				reportURI: "http://www.domain.com",
			},
			expected: `max-age=604800, enforce, report-uri="http://www.domain.com"`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			expectCt, err := ExpectCt(tC.options)
			if err != nil {
				t.Fatal(err)
			}

			check(t, expectCt, "expect-ct", tC.expected)
		})
	}

	t.Run("Fails with negative max age", func(t *testing.T) {
		_, err := ExpectCt(ExpectCtOptions{
			maxAge: -10,
		})
		if err == nil {
			t.Error("Expected to fail")
		}
	})
}
