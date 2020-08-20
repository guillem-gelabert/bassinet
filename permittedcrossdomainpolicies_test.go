package bassinet

import "testing"

func TestPermittedCrossDomainPolicies(t *testing.T) {
	testCases := []struct {
		desc     string
		expected string
		target   int
	}{
		{desc: "with none",
			expected: "none",
			target:   0,
		},
		{desc: "with master-only",
			expected: "master-only",
			target:   1,
		},
		{desc: "with by-content-type",
			expected: "by-content-type",
			target:   2,
		},
		{desc: "with all",
			expected: "all",
			target:   3,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			permittedCrossDomainPolicies, err := PermittedCrossDomainPolicies(tC.target)
			if err != nil {
				t.Fatal(err)
			}

			check(t, permittedCrossDomainPolicies, "X-Permitted-Cross-Domain-Policies", tC.expected)
		})
	}

	t.Run("with invalid policy", func(t *testing.T) {
		_, err := PermittedCrossDomainPolicies(20)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})
}
