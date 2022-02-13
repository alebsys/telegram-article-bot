package devto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultTag       string = ""
	defaultFreshness string = "10"
	defaultLimit     int    = 10
	url                     = "https://dev.to/api/articles"
	dotSymbol               = 9865 // unicode symbol of dot '⚉' https://unicodeplus.com/U+2689
)

var (
	rgxp = [4]string{`^/article\s{1}[a-zA-z]+\s[1-9][0-9]*\s[1-9][0-9]*$`, `^/article\s{1}[a-zA-z]+\s[1-9][0-9]*$`, `^/article\s{1}[a-zA-z]*$`, `^/article$`}
)

type Query struct {
	Tag       string
	Freshness string
	Limit     int
}

type Article struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Score int    `json:"positive_reactions_count"`
}
type Articles []Article

type QueryOption func(*Query) error

// WithTag adds tag to Query or set default value.
func WithTag(tag string) QueryOption {
	return func(q *Query) error {
		q.Tag = defaultTag
		if len(tag) > 0 {
			q.Tag = tag
		}
		return nil
	}
}

// WithFreshness adds freshness to Query or set default value.
func WithFreshness(freshness string) QueryOption {
	return func(q *Query) error {
		q.Freshness = defaultFreshness
		if len(freshness) > 0 {
			q.Freshness = freshness
		}
		return nil
	}
}

// WithLimit adds limit to a Query or set default value.
func WithLimit(limit string) QueryOption {
	return func(q *Query) (err error) {
		q.Limit = defaultLimit
		if len(limit) > 0 {
			q.Limit, err = strconv.Atoi(limit)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ValidateInput parse input sting from user and return true if input is valid.
// User input must be of the format: '/article go 10 5' or '/article go 10' or '/article go' or '/article'
func ValidateInput(input string) bool {
	for i := range rgxp {
		matched, _ := regexp.MatchString(rgxp[i], input)
		if !matched {
			return false
		}
	}
	return true
}

// ParseInput parse user input string and construct Query.
func ParseInput(input string) (*Query, error) {
	args := make([]string, 4, 4)
	argsSplit := strings.Split(input, " ")
	copy(args, argsSplit)

	query := NewQuery(
		WithTag(args[1]),
		WithFreshness(args[2]),
		WithLimit(args[3]),
	)
	return query, nil
}

// NewQuery makes query to DEV.TO API from user input
func NewQuery(opts ...QueryOption) *Query {
	query := new(Query)
	// apply the list of options to Query
	for _, opt := range opts {
		opt(query)
	}
	return query
}

// GetArticles makes request to DEV.TO API and return Articles struct
func GetArticles(tag, fresh string) (*Articles, error) {
	articles := new(Articles)

	url := fmt.Sprintf("%s?tag=%s&top=%s", url, tag, fresh)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when makes http GET from %s: %v", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error when reads from response body: %v", err)
	}

	if err = json.Unmarshal(body, articles); err != nil {
		return nil, fmt.Errorf("error when unmarshal body: %v", err)
	}
	return articles, nil

}

// WriteArticles makes response to user
func (articles *Articles) WriteArticles(limit int) string {
	buf := new(bytes.Buffer)

	for i, a := range *articles {
		if i >= limit {
			break
		}
		buf.WriteRune(dotSymbol)
		buf.WriteString(fmt.Sprintf(" [%s](%s)\n`  Score: %d`\n\n", a.Title, a.Url, a.Score))

	}
	return buf.String()
}
