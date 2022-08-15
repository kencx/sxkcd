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
	"sync"
	"time"
)

// full endpoint: http://xkcd.com/<comic-num>/info.0.json
var (
	baseUrl  = "http://xkcd.com"
	endPoint = "info.0.json"
)

type Comic struct {
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

// TODO safe cancellation
func GetAllComics(total int) ([]*Comic, error) {

	var wg sync.WaitGroup
	ch := make(chan struct{}, 50) // max number of concurrent requests
	comics := make([]*Comic, total)

	for i := 1; i < total+1; i++ {
		wg.Add(1)
		s := strconv.Itoa(i)

		go func(i int, num string) {
			defer wg.Done()
			ch <- struct{}{} // acquire token
			comic, err := GetComic(num)
			if err != nil {
				log.Println(err)
				return
			}
			comics[i-1] = comic
			<-ch // release token
		}(i, s)
	}

	wg.Wait()
	return comics, nil
}

func GetComic(num string) (*Comic, error) {
	url, err := parseEndpoint(baseUrl, num)
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

	var comic Comic
	err = json.Unmarshal([]byte(body), &comic)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("failed to unmarshal due to syntax error at byte offset %d", e.Offset)
		}
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	return &comic, nil
}

func parseEndpoint(baseUrl, num string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint: %v", err)
	}

	u.Path = path.Join(u.Path, num, endPoint)
	return u.String(), nil
}

func main() {
	start := time.Now()
	comics, err := GetAllComics(2000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken: %v\n", time.Since(start))
	fmt.Println(len(comics))
}
