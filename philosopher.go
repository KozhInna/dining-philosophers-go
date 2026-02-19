package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	done		atomic.Int32
}

func (philo *Philosopher) run(ctx context.Context, conf *Config, cancel context.CancelFunc) error {
	//Initial delay
	if err := philo.initialDelay(ctx, conf); err != nil {
		return err
	}

    for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
        // Think
        philo.printAction(philo.id, "is thinking", conf)

		//Eat (take forks, eat, realise forks)
		doneMeal, err := philo.eat(ctx, conf) 
		if err != nil {
			return err
		} 
		if doneMeal {
			if conf.NumPhilosDone.Load() == int32(conf.NumPhilos) {
				cancel()
			}
			return nil
		}

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
	if philo.id % 2 == 0 || conf.NumPhilos == 1  {
		return nil
	}

	delay := conf.TimeToEat
	if philo.id == conf.NumPhilos && conf.NumPhilos != 1 {
		delay = conf.TimeToEat * 2
	}

	return philo.waitOrCancel(ctx, delay)
}


func (philo *Philosopher) waitOrCancel(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return nil
	case <-time.After(duration):
		return nil
	}
}

func (philo *Philosopher) printAction(id int, action string, conf *Config) {
    timestamp := time.Since(conf.StartTime).Milliseconds()
    fmt.Printf("%d %d %s\n", timestamp, id, action)
}

func (philo *Philosopher) eat(ctx context.Context, conf *Config) (bool, error) {

	 // Take forks
	 if err := philo.takeForks(ctx, conf); err != nil {
		return false, err
	}

	// Update last meal time
	philo.mtx.Lock()
	philo.lastMeal = time.Now()
	philo.mtx.Unlock()

	//Eat
	philo.printAction(philo.id, "is eating", conf)
	if err := philo.waitOrCancel(ctx, conf.TimeToEat); err != nil {
		philo.releaseForks()
		return false, err
	}

	// Release forks
	philo.releaseForks()

	// Increment counter and check if done
	philo.timesEaten++
	if conf.TimesToEat > 0 && philo.timesEaten >= conf.TimesToEat {
		philo.done.Store(1)
		conf.NumPhilosDone.Add(1)
		return true, nil  
	}

	return false, nil
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

func (philo *Philosopher) releaseForks() {
	philo.leftFork <- true
	philo.rightFork <- true 
}