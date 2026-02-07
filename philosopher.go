package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Config struct {
    NumPhilos   int
    TimeToDie   time.Duration
    TimeToEat   time.Duration
    TimeToSleep time.Duration
    TimesToEat  int
	StartTime   time.Time 
}

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

var (
	ErrInvalidArgs  = errors.New("invalid number of arguments")
    ErrInvalidValue = errors.New("invalid argument value")
    ErrStarvation   = errors.New("philosopher starved")
)

func parseArgs(args []string) (*Config, error) {
	argc := len(args)
	if argc < 5 || argc > 6 {
		return nil, fmt.Errorf("%w: expected 4-5 arguments, got %d", 
			ErrInvalidArgs, argc - 1)
	}
	numPhilos, err := parsePositiveInt(args[1], "number_of_philosophers")
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
    }
    
    timeToDie, err := parsePositiveInt(args[2], "time_to_die")
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
    }
    
    timeToEat, err := parsePositiveInt(args[3], "time_to_eat")
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
    }
    
    timeToSleep, err := parsePositiveInt(args[4], "time_to_sleep")
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
    }
    
    timesToEat := -1
    if argc == 6 {
        timesToEat, err = parsePositiveInt(args[5], "number_of_times_each_philosopher_must_eat")
        if err != nil {
            return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
        }
    }

	config := &Config{
		NumPhilos:   numPhilos,
        TimeToDie:   time.Duration(timeToDie) * time.Millisecond,
        TimeToEat:   time.Duration(timeToEat) * time.Millisecond,
        TimeToSleep: time.Duration(timeToSleep) * time.Millisecond,
        TimesToEat:  timesToEat,
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidValue, err)
	}

	return config, nil
}

func parsePositiveInt(arg, name string) (int, error) {
    n, err := strconv.Atoi(arg)
    if err != nil {
        return 0, fmt.Errorf("%s must be a valid integer", name)
    }
    if n < 0 {
        return 0, fmt.Errorf("%s must be non-negative", name)
    }
    return n, nil
}

func (conf *Config) validate() error {
    if conf.NumPhilos < 1 {
        return fmt.Errorf("number_of_philosophers must be at least 1")
    }
	if conf.TimeToDie == 0 || conf.TimeToEat == 0 || conf.TimeToSleep == 0 {
		return fmt.Errorf("time must be greater than 0 ms")
	}
	if conf.TimesToEat == 0 {
		return fmt.Errorf("times to eat must be more than 0")
	}
	return nil
}

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
		return fmt.Errorf("%w: philosopher %d died", ErrStarvation, philo.id)
	}
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

func (philo *Philosopher) waitOrCancel(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
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

func (philo *Philosopher) printAction(id int, action string, conf *Config) {
    timestamp := time.Since(conf.StartTime).Milliseconds()
    fmt.Printf("%d %d %s\n", timestamp, id, action)
}