package main

import "testing"

func TestTruncateHostname(t *testing.T) {
	tests := []struct {
		Hostname string
		Expected string
	}{
		{
			Hostname: "bilbo",
			Expected: "bilbo",
		},
		{
			Hostname: "bilbo.shire.lan",
			Expected: "bilbo",
		},
		{
			Hostname: "bilbo.lan",
			Expected: "bilbo",
		},
	}

	for _, tt := range tests {
		got := truncateHostname(tt.Hostname)
		if got != tt.Expected {
			t.Errorf("expected %q, got %q", tt.Expected, got)
		}
	}
}
