package util

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		// skip retry when canceled
		if errors.Is(err, context.Canceled) {
			return err
		}

		if attempts--; attempts > 0 {
			log.Printf("retrying due to err: %v", err)

			jitter := time.Duration((rand.Int63n(int64(sleep))))
			sleep += jitter / 2

			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}
