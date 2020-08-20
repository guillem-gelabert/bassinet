package bassinet

import "testing"

func TestReferrerPolicy(t *testing.T) {
	testCases := []struct {
		desc     string
		target   []int
		expected string
	}{
		{
			desc:     "with same-origin",
			target:   []int{PolicySameOrigin},
			expected: "same-origin,no-referrer",
		},
		{
			desc:     "with multiple policies",
			target:   []int{PolicyOrigin, PolicyUnsafeURL},
			expected: "origin,unsafe-url,no-referrer",
		},
		{
			desc:     "with no-referrer",
			target:   []int{PolicyNoReferrer, PolicyNoReferrer, PolicyNoReferrer},
			expected: "no-referrer",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			referrerPolicy, err := ReferrerPolicy(tC.target)
			if err != nil {
				t.Fatal(err)
			}

			check(t, referrerPolicy, "Referrer-Policy", tC.expected)
		})
	}

	t.Run("with invalid policy", func(t *testing.T) {
		_, err := ReferrerPolicy([]int{123})
		if err == nil {
			t.Error("Expected to fail")
		}
	})
}
