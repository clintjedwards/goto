package models

import "testing"

func TestIsFormattedLink(t *testing.T) {
	tests := map[string]struct {
		url      string
		expected bool
	}{
		"standard": {
			url:      "https://github.com/clintjedwards",
			expected: false,
		},
		"formatted": {
			url:      "https://github.com/clintjedwards/{}",
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(_ *testing.T) {
			if tc.expected != isFormattedLink(tc.url) {
				t.Errorf("isFormatted mismatch for url %q; expected %v, got %v",
					tc.url, tc.expected, isFormattedLink(tc.url))
			}
		})
	}
}

func TestCheckValidID(t *testing.T) {
	tests := map[string]struct {
		id          string
		shouldError bool
	}{
		"simple": {
			id:          "jelly",
			shouldError: false,
		},
		"invalid": {
			id:          "w#eow",
			shouldError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(_ *testing.T) {
			err := checkValidID(tc.id)
			if err != nil && !tc.shouldError {
				t.Errorf("ID %q is invalid when it should be valid", tc.id)
			} else if err == nil && tc.shouldError {
				t.Errorf("ID %q is valid when it should be invalid", tc.id)
			}
		})
	}
}
