package data

import (
	"fmt"
	"log"
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

func ParseAllComics(x map[int]*XkcdComic, e map[int]*ExplainXkcd) (map[int]*Comic, error) {
	comics := make(map[int]*Comic)

	for i := 1; i < len(x)+1; i++ {
		xkcd, ok := x[i]
		if !ok {
			log.Printf("xkcd %d not found\n", i)
			continue
		}
		explain, ok := e[i]
		if !ok {
			log.Printf("explainxkcd %d not found\n", i)
			continue
		}

		c, err := NewComic(*xkcd, *explain)
		if err != nil {
			return nil, err
		}
		comics[i] = c
	}
	return comics, nil
}
