package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
)

// full endpoint: http://xkcd.com/<comic-num>/info.0.json
const (
	XkcdBaseUrl        = "https://xkcd.com"
	xkcdEndpoint       = "info.0.json"
	ExplainBaseUrl     = "https://www.explainxkcd.com/wiki/api.php"
	defaultTimeOut     = 20
	defaultMaxBodySize = 15 * 1024 * 1024
)

type Client struct {
	xkcdUrl     string
	explainUrl  string
	TimeOut     int
	MaxBodySize int64
}

func NewClient(xkcdUrl, explainUrl string) *Client {
	return &Client{
		xkcdUrl:     xkcdUrl,
		explainUrl:  explainUrl,
		TimeOut:     defaultTimeOut,
		MaxBodySize: int64(defaultMaxBodySize),
	}
}

// Retrieves given comic number.
// If 0 is passed, retrieves latest comic.
func (c *Client) RetrieveXkcd(number int) (*XkcdComic, error) {
	url, err := c.parseXkcdEndpoint(number)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %v failed: %v", c.xkcdUrl, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	defer resp.Body.Close()

	var comic XkcdComic
	err = json.Unmarshal([]byte(body), &comic)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("failed to unmarshal %d due to syntax error at byte offset %d", number, e.Offset)
		}
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	return &comic, nil
}

func (c *Client) RetrieveExplain(number int) (*ExplainXkcd, error) {
	url, err := c.parseExplainEndpoint(number)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %v failed: %v", c.xkcdUrl, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	defer resp.Body.Close()

	response := struct {
		Parse struct {
			Wikitext map[string]string
		}
	}{}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("failed to unmarshal %d due to syntax error at byte offset %d", number, e.Offset)
		}
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	wikitext := response.Parse.Wikitext["*"]
	explanation := extractExplanation(wikitext)

	e := &ExplainXkcd{
		Explanation: explanation,
	}
	return e, nil
}

// TODO remove incomplete tag
func extractExplanation(wikitext string) string {
	s := strings.Split(wikitext, "==")
	result := strings.TrimSpace(s[2])
	result = strings.ToValidUTF8(result, " ")
	return result
}

// TODO use errgroup
// Retrieve all xkcd comics up to latest number.
func (c *Client) RetrieveAllXkcds(latest int) (map[int]*XkcdComic, error) {

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	tokens := make(chan struct{}, 50) // max number of concurrent requests
	comics := make(map[int]*XkcdComic, latest)

	for i := 1; i < latest+1; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			defer func() { <-tokens }()

			tokens <- struct{}{}

			if i == 404 {
				return
			}

			comic, err := c.RetrieveXkcd(i)
			if err != nil {
				log.Println(err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			comics[i] = comic
		}(i)
	}
	wg.Wait()
	return comics, nil
}

func (c *Client) RetrieveAllExplain(latest int) (map[int]*ExplainXkcd, error) {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	tokens := make(chan struct{}, 50) // max number of concurrent requests
	comics := make(map[int]*ExplainXkcd, latest)

	for i := 1; i < latest+1; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			defer func() { <-tokens }()

			tokens <- struct{}{}

			if i == 404 {
				return
			}

			comic, err := c.RetrieveExplain(i)
			if err != nil {
				log.Println(err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			comics[i] = comic
		}(i)
	}
	wg.Wait()
	return comics, nil
}

func (c *Client) parseXkcdEndpoint(number int) (string, error) {
	u, err := url.Parse(c.xkcdUrl)
	if err != nil {
		return "", fmt.Errorf("invalid url: %v", err)
	}

	const endPoint = "info.0.json"
	if number == 0 {
		u.Path = path.Join(u.Path, endPoint)
		return u.String(), nil
	}

	n := strconv.Itoa(number)

	u.Path = path.Join(u.Path, n, endPoint)
	return u.String(), nil
}

func (c *Client) parseExplainEndpoint(number int) (string, error) {
	u, err := url.Parse(c.explainUrl)
	if err != nil {
		return "", fmt.Errorf("invalid url: %v", err)
	}

	q := u.Query()
	q.Set("action", "parse")
	q.Set("format", "json")
	q.Set("redirects", "true")
	q.Set("page", strconv.Itoa(number))
	q.Set("prop", "wikitext")
	q.Set("sectiontitle", "Explanation") // only get explanation

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func WriteToFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
