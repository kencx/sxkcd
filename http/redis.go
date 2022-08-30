package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kencx/rkcd/data"
)

// Result is identical to data.Comic but excludes the unnecessary
// transcript and explain attributes that are not rendered but
// usually very large
type Result struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Number int    `json:"num"`
	Alt    string `json:"alt,omitempty"`
	ImgUrl string `json:"img_url"`
	Date   int64  `json:"date"`
}

func comicToResult(i int, c *data.Comic) *Result {
	return &Result{
		Id:     i,
		Title:  c.Title,
		Number: c.Number,
		Alt:    c.Alt,
		ImgUrl: c.ImgUrl,
		Date:   c.Date,
	}
}

func (s *Server) Index() error {
	s.rdb.Do(s.ctx,
		"FT.CREATE", "comics", "ON", "JSON", "PREFIX", "1", "comic:",
		"SCHEMA",
		"$.title", "AS", "title", "TEXT", "WEIGHT", "50",
		"$.alt", "AS", "alt", "TEXT", "WEIGHT", "10",
		"$.transcript", "AS", "transcript", "TEXT", "WEIGHT", "5",
		"$.explanation", "AS", "explanation", "TEXT", "WEIGHT", "1",
		"$.num", "AS", "num", "NUMERIC",
		"$.date", "AS", "date", "NUMERIC",
	)

	pipe := s.rdb.Pipeline()
	for i, c := range s.comics {
		j, err := json.Marshal(&c)
		if err != nil {
			return fmt.Errorf("failed to marshal comic %d: %v", c.Number, err)
		}

		id := strconv.Itoa(i)
		pipe.Do(s.ctx, "JSON.SET", "comic:"+id, "$", j)
	}

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to index: %v", err)
	}
	return nil
}

// returns slice of 100 results
func (s *Server) Search(query string) (int64, []*Result, error) {
	values, err := s.rdb.Do(s.ctx,
		"FT.SEARCH", "comics", query,
		"RETURN", "0",
		"LIMIT", 0, 100,
	).Slice()
	if err != nil {
		return 0, nil, fmt.Errorf("search query failed: %v", err)
	}

	count := values[0].(int64)

	var results []*Result
	for i, v := range values[1:] {
		r := strings.TrimPrefix(v.(string), "comic:")
		id, err := strconv.Atoi(r)
		if err != nil {
			return 0, nil, err
		}
		f := s.comics[id]
		results = append(results, comicToResult(i, f))
	}

	return count, results, nil
}
