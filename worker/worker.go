package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kencx/sxkcd/data"
	"github.com/kencx/sxkcd/redis"
)

type Worker struct {
	rds    *redis.Client
	ticker *time.Ticker
	stop   chan (bool)
	busy   bool
}

func New(client *redis.Client) *Worker {
	return &Worker{
		rds:    client,
		ticker: time.NewTicker(24 * time.Hour),
		stop:   make(chan bool),
		busy:   false,
	}
}

func (w *Worker) Start() error {
	log.Printf("starting worker...")

	go func() {
		for {
			select {
			case <-w.stop:
				w.ticker.Stop()
				log.Println("stopping worker...")
				return
			case <-w.ticker.C:
				err := w.fetchComic()
				if err != nil {
					log.Printf("err processing worker task: %v", err)
				}
			}
		}
	}()
	return nil
}

func (w *Worker) Stop() {
	w.stop <- true
}

func (w *Worker) fetchComic() error {
	if w.busy {
		return errors.New("fetching already in progress")
	}

	w.busy = true
	defer func() {
		w.busy = false
	}()

	client, err := data.NewClient(data.XkcdBaseUrl, data.ExplainBaseUrl)
	if err != nil {
		return fmt.Errorf("failed to create http client: %v", err)
	}
	latest, err := client.RetrieveLatest()
	if err != nil {
		return err
	}

	exists, err := w.rds.ComicExists(latest)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("latest comic #%d already exists, skipping...", latest)
		return nil
	}

	log.Printf("fetching latest comic: #%d", latest)
	start := time.Now()
	comic, err := client.RetrieveComic(latest)
	if err != nil {
		return err
	}

	c, err := json.Marshal(&comic)
	if err != nil {
		return fmt.Errorf("failed to marshal comic: %v", err)
	}
	if err = w.rds.Add(comic.Number, c); err != nil {
		return err
	}

	log.Printf("Successfully fetched comic #%d in %v\n", latest, time.Since(start))
	return nil
}
