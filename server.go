package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Server struct {
	ctx context.Context
	rdb redis.Client
}

func New() *Server {
	return &Server{
		ctx: context.Background(),
		rdb: *redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (s *Server) Run() {

}

func (s *Server) Index(comics []*Comic) error {
	s.rdb.Do(s.ctx,
		"FT.CREATE", "comics", "ON", "JSON", "PREFIX", "1", "comic:",
		"SCHEMA",
		"$.title", "AS", "title", "TEXT", "WEIGHT", "3",
		"$.num", "AS", "number", "NUMERIC",
		"$.alt", "AS", "alt", "TEXT", "WEIGHT", "2",
		"$.transcript", "AS", "transcript", "TEXT",
		"$.img_url", "AS", "url", "TEXT",
		"$.date", "AS", "date", "TEXT",
	)

	pipe := s.rdb.Pipeline()
	for i, c := range comics {
		// TODO replace
		if i == 403 {
			continue
		}

		j, err := json.Marshal(&c)
		if err != nil {
			return fmt.Errorf("failed to marshal comic %d: %v", c.Number, err)
		}

		// TODO check this
		if i < len(comics)-1 {
			pipe.Do(s.ctx, "JSON.SET", "comic:"+strconv.Itoa(c.Number), "$", j)
		}
	}

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to add execute: %v", err)
	}
	return nil
}

// returns slice of first 100 IDs
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
