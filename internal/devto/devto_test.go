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

func TestUnpackSliceToString(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		want   map[int]string
		failed bool
	}{
		{"without args", []string{}, nil, false},
		{"with tag", []string{"go"}, map[int]string{0: "go"}, false},
		{"with failed tag", []string{"go"}, map[int]string{0: "rust"}, true},
		{"with tag and freshness", []string{"go", "10"}, map[int]string{0: "go", 1: "10"}, false},
		{"with tag, freshness and limit", []string{"go", "10", "5"}, map[int]string{0: "go", 1: "10", 2: "5"}, false},
	}
	for _, c := range cases {
		var tag, freshness, limit string
		unpackSliceToString(c.args, &tag, &freshness, &limit)
		for k, v := range c.want {
			if c.args[k] != v {
				if !c.failed {
					t.Errorf("unpackSliceToString: %s; got %v; want %v", c.name, tag, v)
				}
			}
		}
	}
}
