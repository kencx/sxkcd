package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kencx/sxkcd/redis"
	"github.com/kencx/sxkcd/worker"
)

type Server struct {
	rds     *redis.Client
	worker  *worker.Worker
	Static  embed.FS
	Version string
}

func NewServer(uri, version string, static embed.FS) (*Server, error) {
	client, err := redis.New(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %v", err)
	}

	return &Server{
		rds:     client,
		worker:  worker.New(client),
		Version: version,
		Static:  static,
	}, nil
}

func (s *Server) Initialize(filename string, reindex bool) error {
	if filename == "" {
		return fmt.Errorf("no filename provided")
	}

	comics, err := decodeFile(filename)
	if err != nil {
		return err
	}

	start := time.Now()
	log.Printf("Indexing %d comics\n", len(comics))

	err = s.rds.CreateIndex()
	if err != nil {
		if err.Error() == "Index already exists" {
			if reindex {
				s.rds.Reindex()
			} else {
				return fmt.Errorf("%v, include --reindex to replace data", err)
			}
		} else {
			return err
		}
	}

	// comics are not guaranteed to be in order. This depends entirely on the order in
	// which they are fetched.
	err = s.rds.AddBatch(comics)
	if err != nil {
		return err
	}

	log.Printf("Successfully indexed %d comics in %v\n", len(comics), time.Since(start))
	return nil
}

func (s *Server) Verify() error {
	ok, err := s.rds.CheckIndex()
	if err != nil {
		return err
	}

	count, err := s.rds.Count()
	if err != nil {
		return err
	}
	if count > 0 && ok {
		log.Printf("Found existing index and %d comics", count)
	} else {
		return fmt.Errorf("no index or comics found, please provide a file")
	}
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

	err = s.worker.Start()
	if err != nil {
		log.Printf("worker failed to start: %v", err)
	}

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %s, shutting down...", sig.String())

	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.worker.Stop()
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

	count, results, err := s.rds.Search(query)
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
