package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kencx/rkcd/data"
)

const (
	help = `usage: rkcd-cli [options] [file]

  Options:
    -l, --latest    Get latest
    -n, --num       Get by number
    -a, --all	    Get all and save to file
    -v, --version   Version info
    -h, --help	    Show help
`
)

var version string

func main() {
	var (
		num         int
		latest      bool
		all         string
		showVersion bool
	)

	flag.BoolVar(&showVersion, "v", false, "version info")
	flag.BoolVar(&latest, "l", false, "get latest")
	flag.BoolVar(&latest, "latest", false, "get latest")
	flag.IntVar(&num, "n", 0, "get by number")
	flag.IntVar(&num, "num", 0, "get by number")
	flag.StringVar(&all, "a", "", "get all")
	flag.StringVar(&all, "all", "", "get all")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	c, err := data.NewClient(data.XkcdBaseUrl, data.ExplainBaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	if num > 0 {
		comic, err := c.RetrieveComic(num)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(comic)
		os.Exit(0)
	}

	if latest {
		comic, err := c.RetrieveLatest()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(comic)
		os.Exit(0)
	}

	if all != "" {
		latestComicNum, err := c.RetrieveLatest()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Retrieving %d comics from API", latestComicNum-1)
		comics, err := c.RetrieveAllComics(latestComicNum)
		if err != nil {
			log.Fatal(err)
		}

		s, err := json.Marshal(comics)
		if err != nil {
			log.Fatalf("failed to marshal comics: %v", err)
		}

		if err := data.WriteToFile(all, s); err != nil {
			log.Fatalf("failed to write to %s: %v", all, err)
		}
	}
}
