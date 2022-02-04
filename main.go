package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

type Comic struct {
	Title      string
	Number     int `json:"num"`
	Alt        string
	Img        string `json:"img"`
	Transcript string `json:"transcript,omitempty"`
}

func GetComic(num string) Comic {
	baseUrl := "http://xkcd.com"
	url, err := getAPIEndpoint(baseUrl, num)

	if err != nil {
		log.Fatalf("URL invalid: %v", err)
		os.Exit(1)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http GET failed: %v", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
		os.Exit(1)
	}

	var comic Comic
	if err := json.Unmarshal([]byte(body), &comic); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
		os.Exit(1)
	}
	return comic
}

func GetRandomComic() Comic {
	rand.Seed(time.Now().UnixNano())

	c := GetComic("")
	n := rand.Intn(c.Number)

	return GetComic(strconv.Itoa(n))
}

// Parses API path "http://xkcd.com/<comic-num>/info.0.json"
func getAPIEndpoint(baseUrl, num string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid API endpoint")
	}

	u.Path = path.Join(u.Path, num, "info.0.json")
	return u.String(), nil
}

func main() {
	var random = flag.Bool("random", false, "generate random comic")
	flag.Parse()

	for _, v := range os.Args[1:] {
		if *random {
			fmt.Println(GetRandomComic())
			os.Exit(0)
		}

		_, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			log.Fatalf("Input must be integer")
			os.Exit(1)
		}

		fmt.Println(GetComic(v))
	}
}
