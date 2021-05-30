package main

import "testing"

func TestGenerateFormattedLink(t *testing.T) {
	tests := map[string]struct {
		input string
		link  string
		want  string
	}{
		"simple": {
			input: "test/1",
			link:  "https://github.com/clintjedwards/{}",
			want:  "https://github.com/clintjedwards/1",
		},
		"multiple": {
			input: "test/1/2",
			link:  "https://github.com/clintjedwards/{}/{}",
			want:  "https://github.com/clintjedwards/1/2",
		},
		"too few": {
			input: "test/1",
			link:  "https://github.com/clintjedwards/{}/{}",
			want:  "https://github.com/clintjedwards/1/{}",
		},
		"too many": {
			input: "test/1/2/3",
			link:  "https://github.com/clintjedwards/{}/{}",
			want:  "https://github.com/clintjedwards/1/2",
		},
		"with reserved characters": {
			input: "test/?hello",
			link:  "https://github.com/clintjedwards/{}",
			want:  "https://github.com/clintjedwards/?hello",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(_ *testing.T) {
			got := generateFormattedLink(tc.input, tc.link)
			if tc.want != got {
				t.Errorf("malformed generated link; want %q; got %q", tc.want, got)
			}
		})
	}
}
