package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kencx/sxkcd/data"
	"github.com/kencx/sxkcd/http"
)

//go:embed all:ui/build
var static embed.FS

var version string

const (
	help = `usage: sxkcd [server|download] [OPTIONS] [FILE]

  Options:
    -v, --version   Version info
    -h, --help	    Show help

  server:
    -f, --file      Read data from file
    -p, --port      Server port
    -r, --redis     Redis connection URI [host:port]
    -i, --reindex   Reindex existing data with new file

  download:
    -n, --num       Download single comic by number
    -f, --file	    Download all comics to file
`
)

func main() {

	var (
		showVersion bool
		file        string
		port        int
		rds         string
		reindex     bool

		num          int
		downloadFile string
	)

	flag.BoolVar(&showVersion, "v", false, "version info")
	flag.BoolVar(&showVersion, "version", false, "version info")

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.StringVar(&file, "f", "", "read data from file")
	serverCmd.StringVar(&file, "file", "", "read data from file")
	serverCmd.IntVar(&port, "p", 6380, "port")
	serverCmd.IntVar(&port, "port", 6380, "port")
	serverCmd.StringVar(&rds, "r", "localhost:6379", "redis connection URI [host:port]")
	serverCmd.StringVar(&rds, "redis", "localhost:6379", "redis connection URI [host:port]")
	serverCmd.BoolVar(&reindex, "i", false, "reindex with new file")
	serverCmd.BoolVar(&reindex, "reindex", false, "reindex with new file")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadCmd.IntVar(&num, "n", 0, "download comic by number")
	downloadCmd.IntVar(&num, "num", 0, "download comic by number")
	downloadCmd.StringVar(&downloadFile, "f", "", "download all comics to file")
	downloadCmd.StringVar(&downloadFile, "file", "", "download all comics to file")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) <= 1 {
		fmt.Print(help)
		os.Exit(1)
	}

	switch args[0] {
	case "download":
		downloadCmd.Parse(args[1:])

		c := data.NewClient()

		var (
			data []byte
			err  error
		)

		if num < 0 {
			log.Fatalf("comic number must be >= 0")

		} else if num == 0 {
			data, err = c.FetchAll()
			if err != nil {
				log.Fatal(err)
			}

		} else {
			comic, err := c.Fetch(num)
			if err != nil {
				log.Fatal(err)
			}
			data, err = json.Marshal(comic)
			if err != nil {
				log.Fatalf("failed to marshal comic: %v", err)
			}
		}

		if downloadFile != "" {
			if err := os.WriteFile(downloadFile, data, 0644); err != nil {
				log.Fatal(err)
			}

			if num > 0 {
				log.Printf("%d comic(s) downloaded to %s", num, downloadFile)
			} else {
				log.Printf("All comics downloaded to %s", downloadFile)
			}

		} else {
			fmt.Println(string(data))
		}
		os.Exit(0)

	case "server":
		serverCmd.Parse(args[1:])

		if port <= 0 {
			log.Fatalf("Invalid port: %v", port)
		}

		s, err := http.NewServer(rds, version, static)
		if err != nil {
			log.Fatal(err)
		}

		if file != "" {
			if err := s.Initialize(file, reindex); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := s.Verify(); err != nil {
				log.Fatal(err)
			}
		}

		if err := s.Run(port); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Print(help)
		os.Exit(1)
	}
}
