package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestGetXkcdEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, err := NewClient(XkcdBaseUrl, ExplainBaseUrl)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		num := 1
		want := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)

		got, err := c.getXkcdEndpoint(num)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("0", func(t *testing.T) {
		c, err := NewClient(XkcdBaseUrl, ExplainBaseUrl)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		num := 0
		want := "https://xkcd.com/info.0.json"

		got, err := c.getXkcdEndpoint(num)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		c, err := NewClient(XkcdBaseUrl, ExplainBaseUrl)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		num := 5
		for i := 1; i < num+1; i++ {
			want := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)

			got, err := c.getXkcdEndpoint(i)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})
}

func TestGetExplainEndpoint(t *testing.T) {
	c, err := NewClient(XkcdBaseUrl, ExplainBaseUrl)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	num := 5
	for i := 1; i < num+1; i++ {
		want := url.Values{
			"action":       []string{"parse"},
			"format":       []string{"json"},
			"redirects":    []string{"true"},
			"prop":         []string{"wikitext"},
			"sectiontitle": []string{"Explanation"},
			"page":         []string{strconv.Itoa(i)},
		}

		urlString, err := c.getExplainEndpoint(i)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		got, err := url.Parse(urlString)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got.Query(), want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestGetRequestXkcd(t *testing.T) {
	t.Run("xkcd 200", func(t *testing.T) {
		num := 100
		want := XkcdComic{
			Title:  "foo",
			Number: num,
			ImgUrl: "https://example.com",
			Day:    "1",
			Month:  "5",
			Year:   "2016",
		}
		ts, err := testServer(t, num, http.StatusOK, want)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		defer ts.Close()

		c, err := NewClient(ts.URL, "")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		var got XkcdComic
		err = c.getRequest(c.getXkcdEndpoint, num, &got)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("xkcd 404", func(t *testing.T) {
		num := 404
		ts, err := testServer(t, num, http.StatusNotFound, nil)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		defer ts.Close()

		c, err := NewClient(ts.URL, "")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		var got XkcdComic
		err = c.getRequest(c.getXkcdEndpoint, num, &got)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != (XkcdComic{}) {
			t.Errorf("got %v, want %v", got, XkcdComic{})
		}
	})

	t.Run("number < 0", func(t *testing.T) {
		num := -1
		ts, err := testServer(t, num, http.StatusNotFound, nil)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		defer ts.Close()

		c, err := NewClient(ts.URL, "")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		var got XkcdComic
		err = c.getRequest(c.getXkcdEndpoint, num, &got)
		if err == nil {
			t.Fatalf("expected err: number must be >= 0")
		}
	})
}

func TestGetRequestExplain(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		num := 100
		want := ExplainXkcd{
			Explanation: "Hello World",
		}

		ts, err := testServer(t, num, http.StatusOK, want)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		defer ts.Close()

		c, err := NewClient("", ts.URL)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		var got ExplainXkcd
		err = c.getRequest(c.getExplainEndpoint, num, &got)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("404", func(t *testing.T) {
		num := 404

		ts, err := testServer(t, num, http.StatusOK, nil)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		defer ts.Close()

		c, err := NewClient("", ts.URL)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		var got ExplainXkcd
		err = c.getRequest(c.getExplainEndpoint, num, &got)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != (ExplainXkcd{}) {
			t.Errorf("got %v, want %v", got, ExplainXkcd{})
		}
	})

}

func testServer(t *testing.T, number, statusCode int, input interface{}) (*httptest.Server, error) {
	t.Helper()

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(data))
	}))
	return ts, nil
}
