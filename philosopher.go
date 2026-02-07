package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
    id          int
    leftFork    chan bool
	leftForkIndex int
    rightFork   chan bool
	rightForkIndex int
    timesEaten  int
    lastMeal    time.Time
    mtx         sync.Mutex
}

func (philo *Philosopher) run(ctx context.Context, conf *Config) error {
	//Initial delay
	if err := philo.initialDelay(ctx, conf); err != nil {
		return err
	}

    for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
        // Think
        philo.printAction(philo.id, "is thinking", conf)

        // Take forks
		if err := philo.takeForks(ctx, conf); err != nil {
			return err
		}

        // Eat
		philo.mtx.Lock()
		philo.lastMeal = time.Now()
		philo.mtx.Unlock()
		philo.printAction(philo.id, "is eating", conf)
		if err := philo.waitOrCancel(ctx, conf.TimeToEat); err != nil {
			return err
		}

        // Release forks
		philo.leftFork <- true
        philo.rightFork <- true 

        // Sleep
        philo.printAction(philo.id, "is sleeping", conf)
		if err := philo.waitOrCancel(ctx, conf.TimeToSleep); err != nil {
			return err
		}

		// Delay if number of philosophers is odd
		if conf.NumPhilos%2 != 0 {
			if err := philo.waitOrCancel(ctx, 1 * time.Millisecond); err != nil {
				return err
			}
		}
    }
}

func (philo *Philosopher) initialDelay(ctx context.Context, conf *Config) error {
	if philo.id % 2 == 0 {
		return nil
	}

	delay := conf.TimeToEat
	if philo.id == conf.NumPhilos && conf.NumPhilos != 1 {
		delay = conf.TimeToEat * 2
	}

	return philo.waitOrCancel(ctx, delay)
}

func (philo *Philosopher)takeForks(ctx context.Context, conf *Config) error {
	var first, second chan bool

	// Take forks: lower-indexed fork first to prevent deadlock
	if philo.leftForkIndex < philo.rightForkIndex {
		first, second = philo.leftFork, philo.rightFork
	} else  {
		first, second = philo.rightFork, philo.leftFork
	}

	if err := takeFork(ctx, first); err != nil {
		return err
	} 
	philo.printAction(philo.id, "has taken a fork", conf)

	if err := takeFork(ctx, second); err != nil {
		return err
	} 
	philo.printAction(philo.id, "has taken a fork", conf)

	return nil
}

func takeFork(ctx context.Context, fork chan bool) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-fork:
		return nil
	}
}

func (philo *Philosopher) waitOrCancel(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

func (philo *Philosopher) printAction(id int, action string, conf *Config) {
    timestamp := time.Since(conf.StartTime).Milliseconds()
    fmt.Printf("%d %d %s\n", timestamp, id, action)
}