package data

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

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
