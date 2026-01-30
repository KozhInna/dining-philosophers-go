# Dining Philosophers in Go

My solution to the dining philosophers problem implemented in Go using goroutines and channels.

## About

A classic concurrency problem where five philosophers sit at a round table and share five forks. Each philosopher needs both adjacent forks to eat. This creates potential for deadlock.

## Features

- Goroutines for concurrent philosophers
- Channels for fork synchronization
- Even/odd strategy to prevent deadlock

## Usage

```bash
   go run . <num_philos> <time_to_die> <time_to_eat> <time_to_sleep>
```
   
## Example

```bash
   go run . 5 800 50 50
```
   
## Implementation

- Each philosopher runs in its own goroutine
- Forks are represented as buffered channels
- Even-numbered philosophers take left fork first
- Odd-numbered philosophers take right fork first
   
## TODO

- [ ] Death detection monitoring
- [ ] Times must eat parameter
- [ ] Context-based cancellation

# Testing Notes

## Bug Found: Single Philosopher Deadlock

**Discovery:** Manual testing with `./philo 1 800 50 50` revealed a deadlock bug.

**Origin:** With 1 philosopher, leftFork and rightFork point to the same channel. The code attempts to take from the same channel twice, causing a deadlock.

**Status:** Known limitation - documented but not fixed yet. Single philosopher case excluded from automated tests.

## Automated Tests Added

- **Classic 5 philosophers:** Verifies even/odd fork-grabbing strategy
- **Many philosophers (10):** Stress test for scalability

Tests run simulation for 10 seconds to verify no deadlock.

## Race Detection

Running `go test -race` shows no data races - channel-based synchronization is thread-safe.

## Future Work

Add special handling for single philosopher case where leftFork == rightFork.
