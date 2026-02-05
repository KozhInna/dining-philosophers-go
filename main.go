package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	config, err := parseArgs(os.Args)
	if err != nil {
		if errors.Is(err, ErrInvalidArgs) {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Usage: philo <num_philos> <time_to_die> <time_to_eat> <time_to_sleep> [times_must_eat]")
			os.Exit(1)
		}

		if errors.Is(err, ErrInvalidValue) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	
	if err := runSimulation(config); err != nil {
		fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
	}
}