package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kencx/sxkcd/util"
	"golang.org/x/sync/errgroup"
)

const (
	XkcdBaseUrl        = "https://xkcd.com"
	ExplainBaseUrl     = "https://www.explainxkcd.com/wiki/api.php"
	xkcdEndpoint       = "info.0.json"
	defaultTimeOut     = 30
	defaultMaxBodySize = 15 * 1024 * 1024
)

var (
	syntaxErr           *json.SyntaxError
	invalidUnmarshalErr *json.InvalidUnmarshalError
)

type Client struct {
	client      *http.Client
	xkcd        url.URL
	explain     url.URL
	timeOut     int
	maxBodySize int64
}

func NewClient(xkcdUrl, explainUrl string) (*Client, error) {
	xu, err := url.Parse(xkcdUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", xkcdUrl, err)
	}

	eu, err := url.Parse(explainUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", explainUrl, err)
	}

	// explainwiki query values
	values := map[string]string{
		"action":       "parse",
		"format":       "json",
		"redirects":    "true",
		"prop":         "wikitext",
		"sectiontitle": "Explanation",
	}

	q := eu.Query()
	for k, v := range values {
		q.Set(k, v)
	}
	eu.RawQuery = q.Encode()

	c := &Client{
		client: &http.Client{
			Timeout: defaultTimeOut * time.Second,
		},
		maxBodySize: int64(defaultMaxBodySize),
		xkcd:        *xu,
		explain:     *eu,
	}

	return c, nil
}

// Fetch comic by given number
// If 0 is passed, retrieves latest comic.
func (c *Client) Fetch(num int) (*Comic, error) {
	if num < 0 {
		return nil, fmt.Errorf("id cannot be < 0")
	}

	var xkcd Xkcd
	err := c.getWithRetry(c.getXkcdEndpoint, num, &xkcd)
	if err != nil {
		return nil, fmt.Errorf("failed to get xkcd %d: %v", num, err)
	}

	explainWiki := struct {
		Parse struct {
			Wikitext map[string]string
		}
	}{}
	err = c.getWithRetry(c.getExplainEndpoint, num, &explainWiki)
	if err != nil {
		return nil, fmt.Errorf("failed to get explain %d: %v", num, err)
	}
	explain := ExplainXkcd{
		Explanation: extractExplanation(explainWiki.Parse.Wikitext["*"]),
	}

	comic, err := NewComic(xkcd, explain)
	if err != nil {
		return nil, err
	}
	return comic, nil
}

// Fetch latest comic number
func (c *Client) FetchLatestNum() (int, error) {
	var dest Xkcd
	if err := c.getWithRetry(c.getXkcdEndpoint, 0, &dest); err != nil {
		return -1, fmt.Errorf("failed to get latest comic: %v", err)
	}
	return dest.Number, nil
}

// Fetch latest comic
func (c *Client) FetchLatest() (*Comic, error) {
	num, err := c.FetchLatestNum()
	if err != nil {
		return nil, err
	}

	comic, err := c.Fetch(num)
	if err != nil {
		return nil, err
	}
	return comic, nil
}

// Fetch all comics up to latest concurrently.
// This does not guarantee that comics will be in order.
func (c *Client) FetchAll(latest int) (map[int]*Comic, error) {
	if latest < 0 {
		return nil, fmt.Errorf("id cannot be < 0")
	}

	var mu sync.Mutex
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	g, gtx := errgroup.WithContext(ctx)
	g.SetLimit(60)

	comics := make(map[int]*Comic, latest)
	for i := 1; i < latest+1; i++ {
		id := i

		g.Go(func() error {
			if id == 404 {
				return nil
			}
			comic, err := c.Fetch(id)
			if err != nil {
				log.Println(err)
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			comics[id] = comic

			select {
			case <-gtx.Done():
				return gtx.Err()
			default:
				return nil
			}
		})
	}

	if err := g.Wait(); err == nil || err == context.Canceled {
		return comics, nil
	} else {
		return nil, err
	}
}

func (c *Client) FetchAllToFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("no filename provided")
	}

	num, err := c.FetchLatestNum()
	if err != nil {
		return err
	}

	log.Printf("Retrieving %d comics from API", num-1)
	comics, err := c.FetchAll(num)
	if err != nil {
		return err
	}

	s, err := json.Marshal(comics)
	if err != nil {
		return fmt.Errorf("failed to marshal comics: %v", err)
	}

	if err := WriteToFile(filename, s); err != nil {
		return fmt.Errorf("failed to write to %s: %v", filename, err)
	}

	log.Printf("%d comics downloaded to %s", num-1, filename)
	return nil
}

func (c *Client) getWithRetry(f func(int) (string, error), num int, dest interface{}) error {
	err := util.Retry(3, 30*time.Second, func() error {
		return c.get(f, num, dest)
	})
	if err != nil {
		log.Printf("failed to retry, skipping run")
	}
	return nil
}

// Dynamically fetches data from endpoint f with comic number num
func (c *Client) get(f func(int) (string, error), num int, dest interface{}) error {
	// parses url endpoint
	url, err := f(num)
	if err != nil {
		return err
	}

	resp, err := c.client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("request to %v timed out", url)
		}
		return fmt.Errorf("request to %v failed: %v", url, err)
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

// Parses the endpoint at https://xkcd.com/[number]/info.0.json
// concurrent safe as non-pointer is used
func (c Client) getXkcdEndpoint(number int) (string, error) {
	const endPoint = "info.0.json"

	if number < 0 {
		return "", fmt.Errorf("number must be >= 0")
	} else if number == 0 {
		// latest comic
		c.xkcd.Path = path.Join(c.xkcd.Path, endPoint)
	} else {
		n := strconv.Itoa(number)
		c.xkcd.Path = path.Join(c.xkcd.Path, n, endPoint)
	}
	return c.xkcd.String(), nil
}

// Parses the endpoint at https://www.explainxkcd.com page
// concurrent safe as non-pointer is used
func (c Client) getExplainEndpoint(number int) (string, error) {
	q := c.explain.Query()
	q.Set("page", strconv.Itoa(number))
	c.explain.RawQuery = q.Encode()

	return c.explain.String(), nil
}

func extractExplanation(wikitext string) string {

	// extract only first heading (Explanation)
	headingRx := regexp.MustCompile(`\n==[\w\d\s]+==`)
	result := headingRx.Split(wikitext, -1)[1]

	// remove all http/https URLs
	urlRx := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	result = urlRx.ReplaceAllLiteralString(result, "")

	// remove wikitables
	// tableRx := regexp.MustCompile(`\{\|[^()]*\|\}`)
	tableRx := regexp.MustCompile(`\{\|(?s).*\|\}`)
	result = tableRx.ReplaceAllLiteralString(result, "")

	// remove math
	mathRx := regexp.MustCompile(`:*\<math\>(?s).*\<\/math\>`)
	result = mathRx.ReplaceAllLiteralString(result, "")

	// remove incomplete tag
	incompleteRx := regexp.MustCompile(`\{\{incomplete\|(.*)\}\}`)
	result = incompleteRx.ReplaceAllLiteralString(result, "")

	result = strings.TrimSpace(result)
	result = strings.ToValidUTF8(result, "")
	return result
}
