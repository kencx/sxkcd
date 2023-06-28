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
	syntaxError           *json.SyntaxError
	invalidUnmarshalError *json.InvalidUnmarshalError
)

type Client struct {
	Client      *http.Client
	xkcdUrl     url.URL
	explainUrl  url.URL
	TimeOut     int
	MaxBodySize int64
}

func NewClient(xkcdBaseUrl, explainBaseUrl string) (*Client, error) {
	xu, err := url.Parse(xkcdBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", xkcdBaseUrl, err)
	}

	eu, err := url.Parse(explainBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", explainBaseUrl, err)
	}

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
		Client: &http.Client{
			Timeout: defaultTimeOut * time.Second,
		},
		MaxBodySize: int64(defaultMaxBodySize),
		xkcdUrl:     *xu,
		explainUrl:  *eu,
	}

	return c, nil
}

// Retrieve latest comic number
func (c *Client) RetrieveLatest() (int, error) {
	var dest XkcdComic
	if err := c.getRequest(c.getXkcdEndpoint, 0, &dest); err != nil {
		return -1, fmt.Errorf("failed to get latest comic: %v", err)
	}
	return dest.Number, nil
}

// Retrieves given comic by number
// If 0 is passed, retrieves latest comic.
func (c *Client) RetrieveComic(number int) (*Comic, error) {

	var xcomic XkcdComic
	err := c.getRequest(c.getXkcdEndpoint, number, &xcomic)
	if err != nil {
		return nil, fmt.Errorf("failed to get xkcd %d: %v", number, err)
	}

	explainWiki := struct {
		Parse struct {
			Wikitext map[string]string
		}
	}{}
	err = c.getRequest(c.getExplainEndpoint, number, &explainWiki)
	if err != nil {
		return nil, fmt.Errorf("failed to get explain %d: %v", number, err)
	}
	ecomic := ExplainXkcd{
		Explanation: extractExplanation(explainWiki.Parse.Wikitext["*"]),
	}

	comic, err := NewComic(xcomic, ecomic)
	if err != nil {
		return nil, err
	}
	return comic, nil
}

// Retrieves url endpoint from f, performs GET request and unmarshal to dest
func (c *Client) getRequest(f func(int) (string, error), number int, dest interface{}) error {
	url, err := f(number)
	if err != nil {
		return err
	}

	resp, err := c.Client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("request to %v timed out", url)
		}
		return fmt.Errorf("request to %v failed: %v", url, err)
	}

	// unmarshal
	r := http.MaxBytesReader(nil, resp.Body, c.MaxBodySize)
	defer resp.Body.Close()

	decoder := json.NewDecoder(r)
	err = decoder.Decode(dest)

	if err != nil {
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character %d: %v", syntaxError.Offset, err)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("body is empty")
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", c.MaxBodySize)
		// panic when decoding to non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}

// Retrieve all comics up to latest comic number concurrently.
// This does not guarantee that comics will be in order.
func (c *Client) RetrieveAllComics(latest int) (map[int]*Comic, error) {

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
			comic, err := c.RetrieveComic(id)
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

// Parses the endpoint at https://xkcd.com/[number]/info.0.json
// concurrent safe as non-pointer is used
func (c Client) getXkcdEndpoint(number int) (string, error) {
	const endPoint = "info.0.json"

	if number < 0 {
		return "", fmt.Errorf("number must be >= 0")
	} else if number == 0 {
		// latest comic
		c.xkcdUrl.Path = path.Join(c.xkcdUrl.Path, endPoint)
	} else {
		n := strconv.Itoa(number)
		c.xkcdUrl.Path = path.Join(c.xkcdUrl.Path, n, endPoint)
	}
	return c.xkcdUrl.String(), nil
}

// Parses the endpoint at https://www.explainxkcd.com page
// concurrent safe as non-pointer is used
func (c Client) getExplainEndpoint(number int) (string, error) {
	q := c.explainUrl.Query()
	q.Set("page", strconv.Itoa(number))
	c.explainUrl.RawQuery = q.Encode()

	return c.explainUrl.String(), nil
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
