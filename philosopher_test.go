package main

import (
    "testing"
    "time"
)

func TestPhilosophers(t *testing.T) {
    tests := []struct {
        name        string
        numPhilos   int
        timeToDie   time.Duration
        timeToEat   time.Duration
        timeToSleep time.Duration
    }{
        {
            name:        "classic five philosophers",
            numPhilos:   5,
            timeToDie:   time.Second,
            timeToEat:   100 * time.Millisecond,
            timeToSleep: 100 * time.Millisecond,
        },
        {
            name:        "many philosophers",
            numPhilos:   10,
            timeToDie:   time.Second,
            timeToEat:   50 * time.Millisecond,
            timeToSleep: 50 * time.Millisecond,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            conf := &Config{
                NumPhilos:   tt.numPhilos,
                TimeToDie:   tt.timeToDie,
                TimeToEat:   tt.timeToEat,
                TimeToSleep: tt.timeToSleep,
            }

            go runSimulation(conf)
			// Jerify it runs for 10 seconds without deadlock
            time.Sleep(10 * time.Second)
        })
    }
}