package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kencx/rkcd/data"
	"github.com/kencx/rkcd/http"
)

//go:embed all:ui/build
var static embed.FS

var version string

const (
	help = `usage: rkcd [server|download] [OPTIONS] [FILE]

  Options:
    -v, --version   Version info
    -h, --help	    Show help

  server:
    -f, --file      Read data from file
    -p, --port      Server port

  download:
    -l, --latest    Get latest comic number
    -n, --num       Download single comic by number
    -f, --file	    Download all comics to file
`
)

func main() {

	var (
		showVersion bool
		file        string
		port        int

		latest       bool
		num          int
		downloadFile string
	)

	flag.BoolVar(&showVersion, "v", false, "version info")
	flag.BoolVar(&showVersion, "version", false, "version info")

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.StringVar(&file, "f", "", "read data from file")
	serverCmd.StringVar(&file, "file", "", "read data from file")
	serverCmd.IntVar(&port, "p", 6500, "port")
	serverCmd.IntVar(&port, "port", 6500, "port")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadCmd.BoolVar(&latest, "l", false, "get latest comic number")
	downloadCmd.BoolVar(&latest, "latest", false, "get latest comic number")
	downloadCmd.IntVar(&num, "n", 0, "download comic by number")
	downloadCmd.IntVar(&num, "num", 0, "download comic by number")
	downloadCmd.StringVar(&downloadFile, "f", "data/comics.json", "download all comics to file")
	downloadCmd.StringVar(&downloadFile, "file", "data/comics.json", "download all comics to file")

	flag.Parse()

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 1 {
		switch args[0] {
		case "server":
			serverCmd.Parse(args[1:])

		case "download":
			downloadCmd.Parse(args[1:])

			c, err := data.NewClient(data.XkcdBaseUrl, data.ExplainBaseUrl)
			if err != nil {
				log.Fatal(err)
			}

			if num > 0 {
				comic, err := c.RetrieveComic(num)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Comic #%d: %v", num, comic)
				os.Exit(0)
			}

			if latest {
				comic, err := c.RetrieveLatest()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Latest comic: #%d", comic)
				os.Exit(0)
			}

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

			if err := data.WriteToFile(downloadFile, s); err != nil {
				log.Fatalf("failed to write to %s: %v", downloadFile, err)
			}
			log.Printf("%d comics downloaded to %s", latestComicNum-1, downloadFile)
			os.Exit(0)

		default:
			fmt.Print(help)
			os.Exit(1)
		}
	}

	s := http.NewServer(static)
	if file != "" {
		if err := s.ReadFile(file); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Data file must be provided")
	}

	if err := s.Run(port); err != nil {
		log.Fatal(err)
	}
}
