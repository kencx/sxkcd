package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

type Server struct {
	ctx     context.Context
	rdb     redis.Client
	Static  embed.FS
	Version string
}

func NewServer(uri, version string, static embed.FS) (*Server, error) {
	s := &Server{
		ctx: context.Background(),
		rdb: *redis.NewClient(&redis.Options{
			Addr:         uri,
			DialTimeout:  20 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}),
		Version: version,
		Static:  static,
	}

	if err := s.rdb.Ping(s.ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis database: %v", err)
	}

	return s, nil
}

func (s *Server) ReadFile(filename string) error {
	var rc io.ReadCloser

	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		response, err := http.Get(filename)
		if err != nil {
			return fmt.Errorf("failed to get %s: %v", filename, err)
		}
		defer response.Body.Close()
		rc = response.Body

	} else {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", filename, err)
		}
		defer f.Close()
		rc = f
	}

	// Decode JSON to RawMessage to be directly indexed as a ReJSON document
	dec := json.NewDecoder(rc)
	t, err := dec.Token()
	if err != nil {
		return fmt.Errorf("token err %v", err)
	}
	if t.(json.Delim) != '{' {
		return fmt.Errorf("not json object")
	}

	var comics []json.RawMessage
	for dec.More() {
		_, err = dec.Token()
		if err != nil {
			return fmt.Errorf("key err %v", err)
		}

		var val json.RawMessage
		err = dec.Decode(&val)
		if err != nil {
			return fmt.Errorf("decode err %v", err)
		}
		comics = append(comics, val)
	}

	start := time.Now()
	log.Printf("Starting indexing of %d comics\n", len(comics))

	err = s.Index(comics)
	if err != nil {
		return err
	}
	log.Printf("Successfully indexed %d comics in %v\n", len(comics), time.Since(start))
	return nil
}

func (s *Server) Run(port int) error {

	p := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    p,
		Handler: mux,
	}

	mux.HandleFunc("/search", s.searchHandler)
	mux.HandleFunc("/health", s.healthcheckHandler)

	// embed static files
	dir, err := fs.Sub(s.Static, "ui/build")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(dir)))

	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	log.Printf("Server started at %s", p)

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %s, shutting down...", sig.String())

	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(tc); err != nil {
		log.Fatalf("failed to shut down gracefully: %v", err)
	}

	log.Printf("Application gracefully stopped")
	return nil
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	start := time.Now()

	if query == "" {
		log.Printf("invalid request parameters: %v", r.URL.Query())
		errorResponse(w, fmt.Errorf("query parameters required"))
		return
	}

	query = sanitize(query)
	query = parseNumFilter(query)
	query, err := parseDateFilter(query)
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
		return
	}

	count, results, err := s.Search(query)
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
		return
	}
	timeTaken := time.Since(start)
	log.Printf("Query produced %d results in %.3fms: %s", count, timeTaken.Seconds()*1000, query)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"count":      count,
		"results":    results,
		"query_time": timeTaken.Seconds(),
	}); err != nil {
		errorResponse(w, err)
		return
	}
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"version": s.Version,
	}); err != nil {
		errorResponse(w, err)
		return
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
