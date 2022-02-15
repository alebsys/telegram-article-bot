package podcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	rgxp      = `^/podcast\s{1}[a-zA-z]*$`
	url       = "https://dev.to/api/podcast_episodes?username"
	dotSymbol = 9865 // unicode symbol of dot '⚉' https://unicodeplus.com/U+2689
	count     = 15   // count of podcasts
)

type Query struct {
	Tag string
}

type Podcast struct {
	// PodcastTitle string `json:"podcast.title"`
	Title string `json:"title"`
	Path  string `json:"path"`
	// TODO ImageUrl
}

type Podcasts []Podcast

// ValidateInput parse input string from user and return true if input is valid.
// User input must be of the format: '/podcast gotime'
func ValidateInput(input string) (bool, error) {
	return regexp.MatchString(rgxp, input)
}

// ParseInput parse user input string and construct Query.
func ParseInput(input string) *Query {
	args := make([]string, 2)
	argsSplit := strings.Split(input, " ")
	copy(args, argsSplit)

	var tag string
	unpackSliceToString(args[1:], &tag)

	return &Query{Tag: tag}
}

func unpackSliceToString(slice []string, vars ...*string) {
	for i, s := range slice {
		*vars[i] = s
	}
}

// GetPodcasts makes request to DEV.TO API and return Podcasts struct
func GetPodcasts(tag string) (*Podcasts, error) {
	podcasts := new(Podcasts)

	url := fmt.Sprintf("%s=%s", url, tag)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when makes http GET from %s: %v", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error when reads from response body: %v", err)
	}

	if err = json.Unmarshal(body, podcasts); err != nil {
		return nil, fmt.Errorf("error when unmarshal body: %v", err)
	}
	return podcasts, nil

}

// WriteArticles makes response to user
func (podcasts *Podcasts) WritePodcasts() string {
	buf := new(bytes.Buffer)

	for i, p := range *podcasts {
		if i >= count {
			break
		}
		buf.WriteRune(dotSymbol)
		buf.WriteString(fmt.Sprintf(" [%s](https://dev.to%s)\n", p.Title, p.Path))

	}
	return buf.String()
}
