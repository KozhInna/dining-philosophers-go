package main

import (
	"time"
	"context"
	"golang.org/x/sync/errgroup"
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

	g, ctx := errgroup.WithContext(context.Background())
    // Start all philosophers as goroutines with context cancellation
    for _, philo := range philos {
		p := philo

        g.Go(func() error {
            return p.run(ctx, conf)
        })
    }

	g.Go(func() error {
		return monitor(ctx, philos, conf)
	})

	return g.Wait()
}