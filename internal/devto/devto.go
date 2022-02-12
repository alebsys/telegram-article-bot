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
	tag       = ""
	freshness = "10"
	limit     = 10
	url       = "https://dev.to/api/articles"
	dotSymbol = 9865 // unicode symbol of dot '⚉' https://unicodeplus.com/U+2689
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

// ParseInput parse input from user and return true if input valid
// User input must be of the format: '/article go 10 5' or '/article go 10' or '/article go' or '/article'
func ParseInput(input string) bool {
	for i := range rgxp {
		m, _ := regexp.MatchString(rgxp[i], input)
		if m {
			return m
		}
	}
	return false
}

// NewQuery makes query to DEV.TO API from user input
func NewQuery(input string) *Query {
	query := &Query{
		Tag:       tag,
		Freshness: freshness,
		Limit:     limit,
	}

	msg := strings.Split(input, " ")

	switch len(msg) {
	case 2:
		query.Tag = msg[1]
	case 3:
		query.Tag = msg[1]
		query.Freshness = msg[2]
	case 4:
		query.Tag = msg[1]
		query.Freshness = msg[2]
		limit, _ := strconv.Atoi(msg[3])
		if limit > 30 {
			query.Limit = 30
		} else {
			query.Limit, _ = strconv.Atoi(msg[3])
		}
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
