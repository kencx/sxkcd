package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func decodeFile(filename string) ([]json.RawMessage, error) {
	var rc io.ReadCloser

	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		response, err := http.Get(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s: %v", filename, err)
		}
		defer response.Body.Close()
		rc = response.Body

	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", filename, err)
		}
		defer f.Close()
		rc = f
	}

	dec := json.NewDecoder(rc)

	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("token err: %v", err)
	}
	if t.(json.Delim) != '{' {
		return nil, fmt.Errorf("not json object")
	}

	var comics []json.RawMessage
	for dec.More() {
		_, err = dec.Token()
		if err != nil {
			return nil, fmt.Errorf("key err: %v", err)
		}

		var val json.RawMessage
		err = dec.Decode(&val)
		if err != nil {
			return nil, fmt.Errorf("decode err: %v", err)
		}
		comics = append(comics, val)
	}
	return comics, nil
}
