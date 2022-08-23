package http

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kencx/rkcd/data"
)

var version string

type Server struct {
	ctx    context.Context
	rdb    redis.Client
	comics map[int]*data.Comic
	Static embed.FS
}

func NewServer(static embed.FS) *Server {
	return &Server{
		ctx: context.Background(),
		rdb: *redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		}),
		Static: static,
	}
}

func (s *Server) ReadFile(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", filename, err)
	}

	if err := json.Unmarshal(body, &s.comics); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	start := time.Now()
	log.Printf("Starting indexing of %d comics\n", len(s.comics))

	err = s.Index()
	if err != nil {
		return err
	}
	log.Printf("Successfully indexed %d comics in %v\n", len(s.comics), time.Since(start))
	return nil
}

func (s *Server) Run(port int) error {

	mux := http.NewServeMux()
	mux.HandleFunc("/search", s.searchHandler)
	mux.HandleFunc("/health", s.healthcheckHandler)

	// embed static files
	dir, err := fs.Sub(s.Static, "ui/build")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(dir)))

	// graceful shutdown

	p := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at %s", p)
	if err := http.ListenAndServe(p, mux); err != nil {
		return err
	}
	return nil
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	start := time.Now()

	count, results, err := s.Search(query)
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
		return
	}
	timeTaken := time.Since(start)
	log.Printf("Query produced %d results in %vs: %s", count, timeTaken.Seconds(), query)

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
		"version": version,
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
