package repl

import (
	"slices"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input  string
		expect []string
	}{
		{
			input:  " hello world ",
			expect: []string{"hello", "world"},
		},
		{
			input:  "Hello World ",
			expect: []string{"hello", "world"},
		},
		{
			input:  "hi",
			expect: []string{"hi"},
		},
	}

	for _, tt := range cases {
		got := cleanInput(tt.input)
		if !slices.Equal(got, tt.expect) {
			t.Errorf("got %v want %v", got, tt.expect)
		}
	}
}
