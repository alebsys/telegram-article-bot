package devto

import "testing"

func TestValidateInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"article without args", "/article", true},
		{"article with tag", "/article go", true},
		{"article with tag and freshness", "/article go 10", true},
		{"article with tag, freshness and limit", "/article go 10 5", true},
		{"acticle with extra args", "/article go 10 5 1", false},
		{"blank input", "", false},
		{"mistake command", "/mistake", false},
		{"blank command", "/", false},
	}
	for _, c := range cases {
		got := ValidateInput(c.input)
		if got != c.want {
			t.Errorf("ValidateInput: %s; got %v; want %v", c.name, got, c.want)
		}
	}
}
