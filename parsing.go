package main

import (
	"fmt"
	"strconv"
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