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
			target:   PCDPNone,
		},
		{desc: "with master-only",
			expected: "master-only",
			target:   PCDPMasterOnly,
		},
		{desc: "with by-content-type",
			expected: "by-content-type",
			target:   PCDPByContentType,
		},
		{desc: "with all",
			expected: "all",
			target:   PCDPAll,
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
