package main

import (
	"fmt"
	"time"
	"context"
)

func monitor(ctx context.Context, philos []*Philosopher, conf *Config) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, philo := range philos {
				if err := checkIsAlive(philo, conf); err != nil {
					return err
				}
			}
		}
	}
}

func checkIsAlive(philo *Philosopher, conf *Config) error {
	philo.mtx.Lock()
	defer philo.mtx.Unlock()

	if time.Since(philo.lastMeal) > conf.TimeToDie {
		timestamp := time.Since(conf.StartTime).Milliseconds()
        return fmt.Errorf("%d %d died", timestamp, philo.id)
	}
	return nil
}