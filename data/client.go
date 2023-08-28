package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kencx/sxkcd/util"
)

const (
	defaultTimeOut     = 30
	defaultMaxBodySize = 15 * 1024 * 1024
)

type Client struct {
	client      *http.Client
	maxBodySize int64
	ctx         context.Context
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: defaultTimeOut * time.Second,
		},
		maxBodySize: int64(defaultMaxBodySize),
		ctx:         context.Background(),
	}
}

func (c *Client) getXkcd(num int, dest interface{}) error {
	url := buildXkcdURL(num)
	return c.getWithRetry(url, dest)
}

func (c *Client) getExplain(num int, dest interface{}) error {
	url, err := buildExplainURL(num)
	if err != nil {
		return err
	}

	return c.getWithRetry(url, dest)
}

func (c *Client) getWithRetry(url string, dest interface{}) error {
	err := util.Retry(3, 30*time.Second, func() error {
		return c.get(url, dest)
	})
	if err != nil {
		return fmt.Errorf("failed to retry: %w", err)
	}
	return nil
}

func (c *Client) get(url string, dest interface{}) error {
	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("request to %v timed out", url)
		}
		return fmt.Errorf("request to failed: %w", err)
	}

	// unmarshal json
	r := http.MaxBytesReader(nil, resp.Body, c.maxBodySize)
	defer resp.Body.Close()

	decoder := json.NewDecoder(r)
	err = decoder.Decode(dest)

	if err != nil {
		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON at character %d: %v", syntaxErr.Offset, err)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body is empty")
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", c.maxBodySize)
		// panic when decoding to non-nil pointer
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		default:
			return err
		}
	}
	return nil
}
