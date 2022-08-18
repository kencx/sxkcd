package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetrieveXkcd(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		want := `{"title":"testTitle","num":100,"img":"https://example.com","day":"1","month":"5","year":"2016"}`
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, want)
		}))
		defer ts.Close()

		c := NewClient(ts.URL, "")

		num := 100
		comic, err := c.RetrieveXkcd(num)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		got, err := json.Marshal(comic)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if string(got) != want {
			t.Errorf("got %v, want %v", string(got), want)
		}
	})

	t.Run("404", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		c := NewClient(ts.URL, "")

		num := 404
		_, err := c.RetrieveXkcd(num)
		if err == nil {
			t.Fatalf("expected err: failed to unmarshal %d", num)
		}
	})
}
