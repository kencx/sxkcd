package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kencx/rkcd/data"
)

func (s *Server) Index(comics map[int]*data.Comic) error {
	s.rdb.Do(s.ctx,
		"FT.CREATE", "comics", "ON", "JSON", "PREFIX", "1", "comic:",
		"SCHEMA",
		"$.title", "AS", "title", "TEXT", "WEIGHT", "3",
		"$.alt", "AS", "alt", "TEXT", "WEIGHT", "2",
		"$.transcript", "AS", "transcript", "TEXT",
		"$.explanation", "AS", "explanation", "TEXT", "WEIGHT", "1",
		// "$.date", "AS", "date", "TEXT",
	)

	pipe := s.rdb.Pipeline()
	for i, c := range comics {
		j, err := json.Marshal(&c)
		if err != nil {
			return fmt.Errorf("failed to marshal comic %d: %v", c.Number, err)
		}

		id := strconv.Itoa(i)
		pipe.Do(s.ctx, "JSON.SET", "comic:"+id, "$", j)
	}

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to add execute: %v", err)
	}
	return nil
}

// returns slice of first 100 IDs matching query
func (s *Server) Search(query string) (int64, []string, error) {
	values, err := s.rdb.Do(s.ctx,
		"FT.SEARCH", "comics", query,
		"RETURN", "0",
		"LIMIT", "0", "100",
	).Slice()
	if err != nil {
		return 0, nil, fmt.Errorf("search query failed: %v", err)
	}

	count := values[0].(int64)
	var results []string
	for _, v := range values[1:] {
		results = append(results, v.(string))
	}
	return count, results, nil
}
