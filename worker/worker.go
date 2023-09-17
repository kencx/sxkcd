package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kencx/sxkcd/data"
	"github.com/kencx/sxkcd/redis"
	"github.com/kencx/sxkcd/util"
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
					log.Println(err)

					// sleep duration should not be longer than ticker duration
					// signal interrupt will be blocked during sleep
					err := util.Retry(3, 10*time.Second, w.fetchComic)
					if err != nil {
						log.Printf("worker: failed to retry, skipping run")
					}
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
		return fmt.Errorf("worker: fetching already in progress")
	}

	w.busy = true
	defer func() {
		w.busy = false
	}()

	start := time.Now()
	log.Println("worker: fetching latest comic")

	client := data.NewClient()
	latest, err := client.Fetch(0)
	if err != nil {
		return err
	}

	exists, err := w.rds.ComicExists(latest.Number)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("worker: latest comic #%d already exists, skipping...", latest.Number)
		return nil
	}

	c, err := json.Marshal(&latest)
	if err != nil {
		return fmt.Errorf("worker: failed to marshal comic: %w", err)
	}

	if err = w.rds.Add(latest.Number, c); err != nil {
		return err
	}
	log.Printf("worker: successfully fetched comic #%d in %v\n", latest.Number, time.Since(start))
	return nil
}
