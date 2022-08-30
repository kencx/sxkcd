package http

import (
	"testing"
	"time"
)

func TestHandleNumSearch(t *testing.T) {
	t.Run("single number", func(t *testing.T) {
		want := "@num: [256 256]"

		query := "#256"
		got := handleNumSearch(query)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("number range", func(t *testing.T) {
		want := "@num: [256 1000]"

		query := "#256-1000"
		got := handleNumSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		want := "abc"

		query := "abc"
		got := handleNumSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid input with numbers", func(t *testing.T) {
		want := "123"

		query := "123"
		got := handleNumSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestHandleDateSearch(t *testing.T) {
	t.Run("single date", func(t *testing.T) {
		timeNow = func() time.Time {
			ts, err := time.Parse("2006-01-02", "2022-08-30")
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			return ts
		}
		want := "@date: [1640995200 1661817600]"

		query := "@date: 2022-01-01"
		got := handleDateSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("date range", func(t *testing.T) {
		want := "@date: [1640995200 1651968000]"

		query := "@date: 2022-01-01 2022-05-08"
		got := handleDateSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("date range with hypen", func(t *testing.T) {
		want := "@date: [1640995200 1651968000]"

		query := "@date: 2022-01-01-2022-05-08"
		got := handleDateSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("date range with comma", func(t *testing.T) {
		want := "@date: [1640995200 1651968000]"

		query := "@date: 2022-01-01,2022-05-08"
		got := handleDateSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		want := "abc"

		query := "abc"
		got := handleNumSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid input with date", func(t *testing.T) {
		want := "date 2022-01-05"

		query := "date 2022-01-05"
		got := handleNumSearch(query)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
