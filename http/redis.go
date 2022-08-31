package http

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func (s *Server) Index(documents []json.RawMessage) error {
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
	for i, d := range documents {
		id := strconv.Itoa(i)
		pipe.Do(s.ctx, "JSON.SET", "comic:"+id, "$", string(d))
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
		"LIMIT", 0, 100,
	).Slice()
	if err != nil {
		return 0, nil, fmt.Errorf("search query failed: %v", err)
	}

	count := values[0].(int64)

	var results []*Result
	for i, v := range values[1:] {

		// skip comic:[id]
		if _, ok := v.(string); ok {
			continue
		}

		// ["$", data]
		sl, ok := v.([]interface{})
		if !ok {
			return 0, nil, fmt.Errorf("search result could not be parsed")
		}
		b := []byte(sl[1].(string))

		var res Result
		if err := json.Unmarshal(b, &res); err != nil {
			return 0, nil, fmt.Errorf("search result could not be unmarshaled: %v", err)
		}
		res.Id = i
		results = append(results, &res)
	}
	return count, results, nil
}
