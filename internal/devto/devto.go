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
	freshness = "10"
	limit     = 10
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

// ParseInput parse input from user and return true if inpur
func ParseInput(text string) bool {
	m, _ := regexp.MatchString(`^/article\s{1}[a-zA-z]+\s[1-9][0-9]*\s[1-9][0-9]*$`, text)
	if m {
		return m
	}
	m, _ = regexp.MatchString(`^/article\s{1}[a-zA-z]+\s[1-9][0-9]*$`, text)
	if m {
		return m
	}
	m, _ = regexp.MatchString(`^/article\s{1}[a-zA-z]*$`, text)
	if m {
		return m
	}
	m, _ = regexp.MatchString(`^/article$`, text)
	if m {
		return m
	}
	return false
}

// NewQuery makes query to DEV.TO API from user input
func NewQuery(text string) *Query {
	q := new(Query)

	m := strings.Split(text, " ")

	switch len(m) {
	case 1:
		q.Tag = ""
		q.Freshness = freshness
		q.Limit = limit
	case 2:
		q.Tag = m[1]
		q.Freshness = freshness
		q.Limit = limit
	case 3:
		q.Tag = m[1]
		q.Freshness = m[2]
		q.Limit = limit
	case 4:
		q.Tag = m[1]
		q.Freshness = m[2]

		c, _ := strconv.Atoi(m[3])
		if c > 30 {
			q.Limit = 30
		} else {
			q.Limit, _ = strconv.Atoi(m[3])
		}
	}
	return q
}

// GetArticles makes request to DEV.TO API and return Articles struct
func GetArticles(tag, fr string) (*Articles, error) {
	articles := new(Articles)

	url := fmt.Sprintf("https://dev.to/api/articles?tag=%s&top=%s", tag, fr)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when makes http GET: %v", err)
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
		buf.WriteRune(9865)
		buf.WriteString(fmt.Sprintf(" [%s](%s)\n`  Score: %d`\n\n", a.Title, a.Url, a.Score))

	}
	return buf.String()
}
