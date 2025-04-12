package utils

import (
	"testing"
)

func TestRecursiveDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/admin/find?serach=%253Cscript%253Ealert%25281%2529%253C%252Fscript%253E", "/admin/find?serach=<script>alert(1)</script>"},
	}

	for _, test := range tests {
		output := RecursiveDecode(test.input, 3)
		if output != test.expected {
			t.Errorf("Decode failed:\nInput:    %s\nExpected: %s\nGot:      %s\n", test.input, test.expected, output)
		}
	}
}
