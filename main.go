package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	var (
		query string
		index bool
	)

	flag.StringVar(&query, "q", "", "query")
	flag.StringVar(&query, "query", "", "query")

	flag.BoolVar(&index, "i", false, "re-index data")
	flag.BoolVar(&index, "index", false, "re-index data")
	flag.Parse()

	start := time.Now()
	s := New()

	if index {
		latest, err := GetXkcdComic("")
		if err != nil {
			log.Fatalf("failed to get latest comic: %v", err)
		}

		xkcds, err := GetAllComics(latest.Number)
		if err != nil {
			log.Fatal(err)
		}

		e := make([]*ExplainXkcd, len(xkcds))
		comics, err := ParseAllComics(xkcds, e)
		if err != nil {
			log.Fatal(err)
		}

		err = s.Index(comics)
		if err != nil {
			log.Fatal(err)
		}
	}

	count, results, err := s.Search(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d results:\n", count)
	fmt.Println(results)
	fmt.Printf("Time taken: %v\n", time.Since(start))
}
