package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kencx/rkcd/data"
)

var version string

type Server struct {
	ctx context.Context
	rdb redis.Client

	comics map[int]*data.Comic
	Static embed.FS
}

func NewServer(uri string, static embed.FS) (*Server, error) {
	s := &Server{
		ctx: context.Background(),
		rdb: *redis.NewClient(&redis.Options{
			Addr: uri,
		}),
		Static: static,
	}

	if err := s.rdb.Ping(s.ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis database: %v", err)
	}

	return s, nil
}

func (s *Server) ReadFile(filename string) error {
	var (
		body []byte
		err  error
	)

	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		response, err := http.Get(filename)
		if err != nil {
			return fmt.Errorf("failed to get %s: %v", filename, err)
		}
		defer response.Body.Close()

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}
	} else {
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", filename, err)
		}
	}

	if err := json.Unmarshal(body, &s.comics); err != nil {
		return fmt.Errorf("ReadFile: failed to unmarshal data: %v", err)
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
