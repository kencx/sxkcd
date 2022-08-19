package data

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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
	Explanation string `json:"explanation"`
}

type Comic struct {
	Title       string    `json:"title"`
	Number      int       `json:"num"`
	Alt         string    `json:"alt,omitempty"`
	Transcript  string    `json:"transcript,omitempty"`
	ImgUrl      string    `json:"img_url"`
	Explanation string    `json:"explanation"`
	Date        time.Time `json:"date"`
}

func NewComic(x XkcdComic, e ExplainXkcd) (*Comic, error) {

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
		Title:       x.Title,
		Number:      x.Number,
		Alt:         x.Alt,
		Transcript:  x.Transcript,
		ImgUrl:      x.ImgUrl,
		Explanation: e.Explanation,
		Date:        date,
	}

	return c, nil
}

func WriteToFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
