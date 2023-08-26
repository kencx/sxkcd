package data

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestGetXkcdEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		num := 1
		want := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)

		got, err := buildXkcdURL(num)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("0", func(t *testing.T) {
		num := 0
		want := "https://xkcd.com/info.0.json"

		got, err := buildXkcdURL(num)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		num := 5
		for i := 1; i < num+1; i++ {
			want := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)

			got, err := buildXkcdURL(i)
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

		urlString, err := buildExplainURL(i)
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
