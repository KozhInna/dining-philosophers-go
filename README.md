# Dining Philosophers in Go

My solution to the dining philosophers problem implemented in Go using goroutines and channels.

## About

A classic concurrency problem where five philosophers sit at a round table and share five forks. Each philosopher needs both adjacent forks to eat. This creates potential for deadlock.

## Implementation

- **Goroutines** - each philosopher runs concurrently
- **Channels** - buffered channels represent forks
- **Even/odd strategy** - asymmetric fork picking avoids deadlock

## Usage

```bash
   go run . <num_philos> <time_to_die> <time_to_eat> <time_to_sleep>
```
   
## Example

```bash
   go run . 5 800 50 50
```
   
## Status

✅ Core simulation, deadlock prevention, tests  
🚧 Death detection, one philo edge case, times_must_eat, graceful shutdown

## Testing
```bash
go test -v
go test -race
