package main

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
	"time"
)

// full endpoint: http://xkcd.com/<comic-num>/info.0.json
var (
	baseUrl  = "http://xkcd.com"
	endPoint = "info.0.json"
)

type XkcdComic struct {
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title,omitempty"`
	Number     int    `json:"num"`
	Alt        string `json:"alt,omitempty"`
	ImgUrl     string `json:"img"`
	Transcript string `json:"transcript,omitempty"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
}

type ExplainXkcd struct {
}

type Comic struct {
	Title      string    `json:"title"`
	Number     int       `json:"num"`
	Alt        string    `json:"alt,omitempty"`
	Transcript string    `json:"transcript,omitempty"`
	ImgUrl     string    `json:"img_url"`
	Date       time.Time `json:"date"`
}

func NewComic(x XkcdComic) (*Comic, error) {

	// TODO better way to check leading 0 in day and month
	if len([]rune(x.Day)) <= 1 {
		x.Day = strings.Join([]string{"0", x.Day}, "")
	}

	if len([]rune(x.Month)) <= 1 {
		x.Month = strings.Join([]string{"0", x.Month}, "")
	}

	dateString := strings.Join([]string{x.Year, x.Month, x.Day}, "-")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	c := &Comic{
		Title:      x.Title,
		Number:     x.Number,
		Alt:        x.Alt,
		Transcript: x.Transcript,
		ImgUrl:     x.ImgUrl,
		Date:       date,
	}

	return c, nil
}

func ParseAllComics(x []*XkcdComic, e []*ExplainXkcd) ([]*Comic, error) {
	comics := make([]*Comic, len(x)+1)

	for i := 0; i < len(x); i++ {
		if i == 403 {
			continue
		}
		c, err := NewComic(*x[i])
		if err != nil {
			return nil, err
		}
		comics[i] = c
	}
	return comics, nil
}

// TODO use errgroup
// TODO convert to a dict to skip 404
func GetAllComics(total int) ([]*XkcdComic, error) {
	var wg sync.WaitGroup
	tokens := make(chan struct{}, 50) // max number of concurrent requests
	comics := make([]*XkcdComic, total)

	for i := 1; i < total+1; i++ {
		wg.Add(1)

		go func(i int) {
			defer func() { <-tokens }()
			defer wg.Done()

			tokens <- struct{}{}
			if i == 404 {
				return
			}

			comic, err := GetXkcdComic(strconv.Itoa(i))
			if err != nil {
				log.Println(err)
				return
			}

			// no need for mutex as all goroutines write to different memory
			comics[i-1] = comic
		}(i)
	}
	wg.Wait()
	return comics, nil
}

func GetXkcdComic(num string) (*XkcdComic, error) {
	url, err := parseXkcdEndpoint(baseUrl, num)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %v failed: %v", url, err)
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
			return nil, fmt.Errorf("failed to unmarshal %s due to syntax error at byte offset %d", num, e.Offset)
		}
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	return &comic, nil
}

func parseXkcdEndpoint(baseUrl, num string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %v", err)
	}

	u.Path = path.Join(u.Path, num, endPoint)
	return u.String(), nil
}
