package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	syntaxErr           *json.SyntaxError
	invalidUnmarshalErr *json.InvalidUnmarshalError
)

// Fetch all comics latest concurrently.
// This does not guarantee that comics will be in order.
func (c *Client) FetchAll() ([]byte, error) {
	// fetch latest
	latest, err := c.Fetch(0)
	if err != nil {
		return nil, err
	}
	num := latest.Number

	log.Printf("Retrieving %d comics from API", num-1)

	var mu sync.Mutex
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	g, gtx := errgroup.WithContext(ctx)
	c.ctx = gtx
	g.SetLimit(50)

	// reset context
	defer func() {
		c.ctx = context.Background()
	}()

	progress := 0

	comics := make(map[int]*Comic, num)
	for i := 1; i < num+1; i++ {
		id := i

		g.Go(func() error {
			if id == 404 {
				return nil
			}
			comic, err := c.Fetch(id)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			comics[id] = comic

			progress += 1
			return nil
		})

		if progress%200 == 0 && progress != 0 {
			log.Printf("Downloaded: %d/%d comics\n", progress, num-1)
		}
	}

	err = g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("cancelled due to signal interrupt")
		} else {
			return nil, err
		}
	}

	s, err := json.Marshal(comics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comics: %v", err)
	}
	return s, nil
}

// Fetch comic by given number
// If 0 is passed, retrieves latest comic.
func (c *Client) Fetch(num int) (*Comic, error) {
	if num < 0 {
		return nil, fmt.Errorf("id cannot be < 0")
	}

	var xkcd Xkcd
	err := c.getXkcd(num, &xkcd)
	if err != nil {
		return nil, fmt.Errorf("failed to get xkcd %d: %w", num, err)
	}

	explainWiki := struct {
		Parse struct {
			Wikitext map[string]string
		}
	}{}
	err = c.getExplain(num, &explainWiki)
	if err != nil {
		return nil, fmt.Errorf("failed to get explain %d: %w", num, err)
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

func extractExplanation(wikitext string) string {
	if wikitext == "" {
		return ""
	}

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
