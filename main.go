package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kencx/rkcd/data"
	"github.com/kencx/rkcd/http"
)

func main() {

	var (
		query  string
		index  bool
		output string
		input  string
	)
	flag.BoolVar(&index, "i", false, "re-index data")
	flag.BoolVar(&index, "index", false, "re-index data")

	flag.StringVar(&input, "f", "", "input json data filepath")
	flag.StringVar(&input, "file", "", "input json data filepath")

	flag.StringVar(&output, "w", "", "output to json file")
	flag.StringVar(&output, "write", "", "output to json file")

	flag.StringVar(&query, "q", "", "query")
	flag.StringVar(&query, "query", "", "query")

	flag.Parse()

	start := time.Now()
	s := http.NewServer()

	if index {
		c, err := data.NewClient(data.XkcdBaseUrl, data.ExplainBaseUrl)
		if err != nil {
			log.Fatal(err)
		}

		latest, err := c.RetrieveLatest()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Retrieving %d comics from API", latest-1)
		comics, err := c.RetrieveAllComics(latest)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(comics))

		// if output != "" {
		// 	dataset, err := json.Marshal(comics)
		// 	if err != nil {
		// 		log.Fatalf("failed to marshal dataset: %v", err)
		// 	}
		//
		// 	err = data.WriteToFile(output, dataset)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	log.Printf("Dataset written to %s successfully", output)
		// }

		log.Printf("Starting indexing of %d comics\n", len(comics))
		err = s.Index(comics)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Indexed %d comics successfully\n", len(comics))
	}

	if input != "" {
		dataset, err := ioutil.ReadFile(input)
		if err != nil {
			log.Fatalf("failed to read file %s: %v", input, err)
		}

		var inputData []*data.Comic
		if err := json.Unmarshal(dataset, &inputData); err != nil {
			log.Fatalf("failed to unmarshal dataset: %v", err)
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
