package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kencx/rkcd/http"
)

var version string

func main() {

	var (
		file        string
		port        int
		showVersion bool
	)

	flag.StringVar(&file, "f", "", "read data from file")
	flag.StringVar(&file, "file", "", "read data from file")
	flag.IntVar(&port, "p", 6500, "port")
	flag.IntVar(&port, "port", 6500, "port")
	flag.BoolVar(&showVersion, "v", false, "version info")
	flag.BoolVar(&showVersion, "version", false, "version info")

	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	s := http.NewServer()

	fmt.Println(file)
	if file != "" {
		if err := s.ReadFile(file); err != nil {
			log.Fatal(err)
		}
	}

	if err := s.Run(port); err != nil {
		log.Fatal(err)
	}
}
