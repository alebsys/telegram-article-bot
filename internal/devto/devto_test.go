package devto

import (
	"fmt"
	"testing"
)

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
		got, err := ValidateInput(c.input)
		if err != nil && got != c.want {
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

func TestWithTag(t *testing.T) {
	type ret struct {
		query  *Query
		err    error
		failed bool
	}
	cases := []struct {
		name string
		tag  string
		want ret
	}{
		{"without args, use defaults", "", ret{&Query{}, nil, false}},
		{"with tag", "go", ret{&Query{Tag: "go"}, nil, false}},
		{"with failed tag", "rust", ret{&Query{Tag: "erlang"}, fmt.Errorf("bad tag"), true}},
		{"with error", "perl", ret{&Query{Tag: "perl"}, fmt.Errorf("bad tag"), false}},
	}
	for _, c := range cases {
		query, err := NewQuery(WithTag(c.tag))
		if err != nil && c.want.err == nil && !c.want.failed {
			t.Errorf("WithTag: %s; got error %v; want error %v", c.name, err, c.want.err)
		}
		if c.want.err == nil && query != nil && query.Tag != c.want.query.Tag {
			t.Errorf("WithTag: %s; got %v; want %v", c.name, query.Tag, c.want.query.Tag)
		}
	}
}
