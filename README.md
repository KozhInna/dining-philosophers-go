# Dining Philosophers in Go

My solution to the dining philosophers problem implemented in Go using goroutines and channels.

## About

A classic concurrency problem where five philosophers sit at a round table and share five forks. Each philosopher needs both adjacent forks to eat. This creates potential for deadlock.

## Implementation

### Concurrency Handling
- **Goroutines** - each philosopher runs as an independent goroutine
- **Channels** - buffered channels represent forks
- **Context** - `errgroup.WithContext` for coordinated cancellation
- **Mutex** - protects shared state

### Deadlock Prevention
- **Resource hierarchy** - always acquire lower-indexed fork first to break circular wait
- **Even/odd initial delays** - stagger philosopher start

### Error Handling
- **Sentinel errors** - `ErrInvalidArgs`, `ErrInvalidValue`
- **Error wrapping** - `%w` with context, checked via `errors.Is()`

## Installation

```bash
git clone https://github.com/KozhInna/dining-philosophers-go.git philo
cd philo
```

## Usage

```bash
   go run . <num_philos> <time_to_die> <time_to_eat> <time_to_sleep> [times_must_eat]
```
   
## Example

```bash
go run . 5 800 200 200        # Runs indefinitely
go run . 5 610 200 200        # Tighter timing
go run . 200 410 200 200      # Stress test: 200 philosophers
go run . 5 800 200 200 5      # Stops after each eats 5 times
```

## Testing
```bash
go test -v
go test -race
```