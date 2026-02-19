package main

import (
	"time"
	"context"
	"golang.org/x/sync/errgroup"
	"fmt"
)

func runSimulation(conf *Config) error {
    conf.StartTime = time.Now()

    // Create forks
    forks := make([]chan bool, conf.NumPhilos)
    for i := range forks {
        forks[i] = make(chan bool, 1)
        forks[i] <- true
    }

    // Create philosophers
    philos := make([]*Philosopher, conf.NumPhilos)
    for i := 0; i < conf.NumPhilos; i++ {
        philos[i] = &Philosopher{
            id:         i + 1,
            leftFork:   forks[i],
			leftForkIndex: i,
            rightFork:  forks[(i+1)%conf.NumPhilos],
			rightForkIndex: (i+1)%conf.NumPhilos,
            timesEaten: 0,
            lastMeal:   conf.StartTime,
        }
    }

    // Parent context
    parentCtx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Child context
    g, ctx := errgroup.WithContext(parentCtx)

    // Start all philosophers as goroutines with context cancellation
    for _, philo := range philos {
		p := philo

        g.Go(func() error {
            return p.run(ctx, conf, cancel)
        })
    }

	g.Go(func() error {
		return monitor(ctx, philos, conf)
	})

	err := g.Wait()
	if err == nil {
		timestamp := time.Since(conf.StartTime).Milliseconds()
		fmt.Printf("%d all philosophers have eaten enough\n", timestamp)
	}
	return err
}