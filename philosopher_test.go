package main

import (
	"testing"
	"time"
)

// TestPhilosophers_Survive verifies philosophers complete successfully with good parameters
func TestPhilosophers_Survive(t *testing.T) {
	tests := []struct {
		name        string
		numPhilos   int
		timeToDie   time.Duration
		timeToEat   time.Duration
		timeToSleep time.Duration
		timesToEat  int
	}{
		{
			name:        "classic five philosophers",
			numPhilos:   5,
			timeToDie:   500 * time.Millisecond,
			timeToEat:   100 * time.Millisecond,
			timeToSleep: 100 * time.Millisecond,
			timesToEat:  3,
		},
		{
			name:        "many philosophers",
			numPhilos:   200,
			timeToDie:   110 * time.Millisecond,
			timeToEat:   50 * time.Millisecond,
			timeToSleep: 50 * time.Millisecond,
			timesToEat:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				NumPhilos:   tt.numPhilos,
				TimeToDie:   tt.timeToDie,
				TimeToEat:   tt.timeToEat,
				TimeToSleep: tt.timeToSleep,
				TimesToEat:  tt.timesToEat,
			}

			err := runSimulation(conf)

			t.Logf("Received error: %v", err)
			
			if err != nil {
				t.Errorf("expected clean completion, got: %v", err)
			}
		})
	}
}

// TestPhilosophers_Die verifies philosophers die with impossible timing parameters
func TestPhilosophers_Die(t *testing.T) {
	conf := &Config{
		NumPhilos:   5,
		TimeToDie:   50 * time.Millisecond,
		TimeToEat:   200 * time.Millisecond,
		TimeToSleep: 100 * time.Millisecond,
		TimesToEat:  -1,
	}

	err := runSimulation(conf)

	t.Logf("Received error: %v", err)
	
	if err == nil {
		t.Fatal("expected simulation to fail with impossible parameters, but got nil")
	}
}

// TestSinglePhilosopher verifies single philosopher edge case (cannot eat with one fork)
func TestSinglePhilosopher(t *testing.T) {
	conf := &Config{
		NumPhilos:   1,
		TimeToDie:   200 * time.Millisecond,
		TimeToEat:   100 * time.Millisecond,
		TimeToSleep: 100 * time.Millisecond,
		TimesToEat:  -1,
	}

	err := runSimulation(conf)

	t.Logf("Received error: %v", err)
	
	if err == nil {
		t.Fatal("single philosopher should die (can't get two forks)")
	}
}