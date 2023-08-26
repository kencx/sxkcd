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
	serverCmd.StringVar(&rds, "r", "redis:6379", "redis connection URI [host:port]")
	serverCmd.StringVar(&rds, "redis", "redis:6379", "redis connection URI [host:port]")
	serverCmd.BoolVar(&reindex, "i", false, "reindex with new file")
	serverCmd.BoolVar(&reindex, "reindex", false, "reindex with new file")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadCmd.IntVar(&num, "n", 0, "download comic by number")
	downloadCmd.IntVar(&num, "num", 0, "download comic by number")
	downloadCmd.StringVar(&downloadFile, "f", "data/comics.json", "download all comics to file")
	downloadCmd.StringVar(&downloadFile, "file", "data/comics.json", "download all comics to file")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 1 {
		switch args[0] {
		case "download":
			downloadCmd.Parse(args[1:])

			c := data.NewClient()

			if num > 0 {
				comic, err := c.Fetch(num)
				if err != nil {
					log.Fatal(err)
				}
				b, err := json.Marshal(comic)
				if err != nil {
					log.Fatalf("failed to marshal comic: %v", err)
				}
				fmt.Println(string(b))
				os.Exit(0)
			}

			if err := c.FetchAll(downloadFile); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)

		case "server":
			serverCmd.Parse(args[1:])

			if rds == "" {
				log.Fatal("Redis connection URI must be provided")
			}
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
}
