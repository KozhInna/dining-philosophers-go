# Dining Philosophers in Go

My solution to the dining philosophers problem implemented in Go using goroutines and channels.

## About

A classic concurrency problem where five philosophers sit at a round table and share five forks. Each philosopher needs both adjacent forks to eat. This creates potential for deadlock.

This implementation includes:
- Go goroutines for concurrent execution
- Channels for synchronization
- Deadlock prevention

## Usage
```bash
go run main.go
```