package http

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var timeNow = time.Now

func sanitize(query string) string {
	chars := []string{"{", "}", "[", "]", "(", ")", "~", ";", `"`, `'`, "%"}
	for _, c := range chars {
		query = strings.ReplaceAll(query, c, "")
	}
	return query
}

func parseNumFilter(query string) string {
	rx := regexp.MustCompile(`\#([0-9]+)\-?([0-9]*)`)
	if rx.MatchString(query) {
		matches := rx.FindStringSubmatch(query)
		if matches[2] != "" {
			query = fmt.Sprintf("@num: [%s %s]", matches[1], matches[2])
		} else {
			query = fmt.Sprintf("@num: [%s %s]", matches[1], matches[1])
		}
	}
	return query
}

func parseDateFilter(query string) (string, error) {
	rx := regexp.MustCompile(`@date:\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})\s?[,-]?\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})?`)
	if rx.MatchString(query) {
		matches := rx.FindStringSubmatch(query)
		from, err := epoch(matches[1])
		if err != nil {
			return "", err
		}
		var to int64

		if matches[2] != "" {
			to, err = epoch(matches[2])
			if err != nil {
				return "", err
			}
		} else {
			to = timeNow().Unix()
		}
		query = fmt.Sprintf("@date: [%d %d]", from, to)
	}
	return query, nil
}

// convert datetime string to epoch time
func epoch(s string) (int64, error) {
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return -1, fmt.Errorf("unable to parse datetime %s: %v", s, err)
	}
	return date.Unix(), nil
}
