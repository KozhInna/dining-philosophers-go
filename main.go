package main

import (
	"fmt"
	"os"
)

func main() {
	config, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Printf("Usage: %s <num_philos> <time_to_die> <time_to_eat> <time_to_sleep> [times_must_eat]\n", os.Args[0])
		os.Exit(1)
	}

	if err := config.validate(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
	
	runSimulation(config)
}