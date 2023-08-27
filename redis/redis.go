package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kencx/sxkcd/data"
	"github.com/redis/go-redis/v9"
)

const (
	Index     = "comics"
	KeyPrefix = "comic:"
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

type Client struct {
	ctx context.Context
	rd  *redis.Client
}

func New(uri string) (*Client, error) {
	r := &Client{
		ctx: context.Background(),
		rd: redis.NewClient(&redis.Options{
			Addr:         uri,
			DialTimeout:  20 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}),
	}
	if err := r.rd.Ping(r.ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis database: %v", err)
	}
	return r, nil
}

// Create JSON index with key comic:[num]
func (r *Client) CreateIndex() error {
	return r.rd.Do(r.ctx,
		"FT.CREATE", Index, "ON", "JSON", "PREFIX", "1", KeyPrefix,
		"SCHEMA",
		"$.title", "AS", "title", "TEXT", "WEIGHT", "50",
		"$.alt", "AS", "alt", "TEXT", "WEIGHT", "10",
		"$.transcript", "AS", "transcript", "TEXT", "WEIGHT", "5",
		"$.explanation", "AS", "explanation", "TEXT", "WEIGHT", "1",
		"$.num", "AS", "num", "NUMERIC",
		"$.date", "AS", "date", "NUMERIC",
	).Err()
}

func (r *Client) Reindex() error {
	err := r.rd.FlushDB(r.ctx).Err()
	if err != nil {
		return err
	}

	err = r.CreateIndex()
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) CheckIndex() (bool, error) {
	indexes, err := r.rd.Do(r.ctx, "FT._LIST").Slice()
	if err != nil {
		return false, err
	}

	if len(indexes) > 0 {
		if indexes[0] == Index {
			return true, nil
		}
	}
	return false, nil
}

func (r *Client) Count() (int, error) {
	count, err := r.rd.DBSize(r.ctx).Result()
	if err != nil {
		return -1, err
	}
	return int(count), nil
}

// Add document if not already exists
func (r *Client) Add(id int, comic []byte) error {
	id_str := strconv.Itoa(id - 1)

	exists, err := r.rd.Exists(r.ctx, KeyPrefix+id_str).Result()
	if err != nil {
		if err != redis.Nil {
			return fmt.Errorf("failed to add comic: %v", err)
		}
	}
	if exists != 0 {
		fmt.Printf("comic %v already present", KeyPrefix+id_str)
		return nil
	}

	err = r.rd.Do(r.ctx, "JSON.SET", KeyPrefix+id_str, "$", string(comic)).Err()
	if err != nil {
		return fmt.Errorf("failed to add comic: %v", err)
	}
	return nil
}

func (r *Client) AddBatch(documents []data.Comic) error {
	pipe := r.rd.Pipeline()
	for i, d := range documents {
		id := strconv.Itoa(i)
		j, err := json.Marshal(&d)
		if err != nil {
			return fmt.Errorf("failed to marshal comic %d: %v", d.Number, err)
		}
		pipe.Do(r.ctx, "JSON.SET", KeyPrefix+id, "$", j)
	}

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to index: %v", err)
	}
	return nil
}

// This checks for existing comic with the comic number $.num in the schema.
// The comic number is entirely different from the Redis key "comic:id" due to
// zero indexing and missing comic 404.
func (r *Client) ComicExists(num int) (bool, error) {
	query := fmt.Sprintf("@num: [%d %d]", num, num)
	count, result, err := r.Search(query)
	if err != nil {
		return false, fmt.Errorf("failed to find comic %d: %v", num, err)
	}
	return (count == 1 && result != nil), nil
}

// returns slice of 100 results
func (r *Client) Search(query string) (int64, []*Result, error) {
	values, err := r.rd.Do(r.ctx,
		"FT.SEARCH", Index, query,
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
