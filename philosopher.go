package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
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
    rightFork   chan bool
    timesEaten  int
    lastMeal    time.Time
    mtx         sync.Mutex
}

func parseArgs(args []string) (*Config, error) {
	argc := len(args)
	if argc < 5 || argc > 6 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	numPhilos, err := parsePositiveInt(args[1], "number_of_philosophers")
    if err != nil {
        return nil, err
    }
    
    timeToDie, err := parsePositiveInt(args[2], "time_to_die")
    if err != nil {
        return nil, err
    }
    
    timeToEat, err := parsePositiveInt(args[3], "time_to_eat")
    if err != nil {
        return nil, err
    }
    
    timeToSleep, err := parsePositiveInt(args[4], "time_to_sleep")
    if err != nil {
        return nil, err
    }
    
    timesToEat := -1
    if argc == 6 {
        timesToEat, err = parsePositiveInt(args[5], "number_of_times_each_philosopher_must_eat")
        if err != nil {
            return nil, err
        }
    }
	return &Config{
        NumPhilos:   numPhilos,
        TimeToDie:   time.Duration(timeToDie) * time.Millisecond,
        TimeToEat:   time.Duration(timeToEat) * time.Millisecond,
        TimeToSleep: time.Duration(timeToSleep) * time.Millisecond,
        TimesToEat:  timesToEat,
    }, nil
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

func runSimulation(conf *Config) {
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
            rightFork:  forks[(i+1)%conf.NumPhilos],
            timesEaten: 0,
            lastMeal:   conf.StartTime,
        }
    }

    // Start all philosophers as goroutines
    var wg sync.WaitGroup
    for i, philo := range philos {
        wg.Add(1)
        isEven := i%2 == 0
        
        go func(p *Philosopher, even bool) {
            defer wg.Done()
            p.run(conf, even)
        }(philo, isEven)
    }

    // Wait forever (or until Ctrl+C)
    wg.Wait()
}

func (philo *Philosopher) run(conf *Config, isEven bool) {
    for {
        // Think
        philo.printAction(philo.id, "is thinking", conf)

        // Take forks (even/odd strategy to avoid deadlock)
        if isEven {
            <-philo.leftFork
            philo.printAction(philo.id, "has taken a fork", conf)
            <-philo.rightFork
            philo.printAction(philo.id, "has taken a fork", conf)
        } else {
            <-philo.rightFork
            philo.printAction(philo.id, "has taken a fork", conf)
            <-philo.leftFork
            philo.printAction(philo.id, "has taken a fork", conf)
        }

        // Eat
        philo.printAction(philo.id, "is eating", conf)
        time.Sleep(conf.TimeToEat)

        // Release forks
        philo.leftFork <- true
        philo.rightFork <- true

        // Sleep
        philo.printAction(philo.id, "is sleeping", conf)
        time.Sleep(conf.TimeToSleep)
    }
}

func (philo *Philosopher) printAction(id int, action string, conf *Config) {
    timestamp := time.Since(conf.StartTime).Milliseconds()
    fmt.Printf("%d %d %s\n", timestamp, id, action)
}