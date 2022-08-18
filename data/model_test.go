package data

import (
	"testing"
	"time"
)

func TestNewComic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testXkcdComic := XkcdComic{
			Title:  "foo",
			Number: 250,
			Day:    "25",
			Month:  "11",
			Year:   "2015",
		}
		testExplain := ExplainXkcd{}

		expectedDate, err := time.Parse("2006-01-02", "2015-11-25")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		want := Comic{
			Title:  "foo",
			Number: 250,
			Date:   expectedDate,
		}

		got, err := NewComic(testXkcdComic, testExplain)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !assertComicEqual(*got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("parse date", func(t *testing.T) {
		testXkcdComic := XkcdComic{
			Title:  "foo",
			Number: 250,
			Day:    "1",
			Month:  "1",
			Year:   "2010",
		}
		testExplain := ExplainXkcd{}

		expectedDate, err := time.Parse("2006-01-02", "2010-01-01")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		want := Comic{
			Title:  "foo",
			Number: 250,
			Date:   expectedDate,
		}

		got, err := NewComic(testXkcdComic, testExplain)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !assertComicEqual(*got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func assertComicEqual(got, want Comic) bool {
	return got.Date.Equal(want.Date) &&
		got.Title == want.Title &&
		got.Number == want.Number &&
		got.Alt == want.Alt &&
		got.Transcript == want.Transcript &&
		got.ImgUrl == want.ImgUrl
}
