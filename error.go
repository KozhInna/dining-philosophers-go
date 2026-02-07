package main

import (
	"errors"
)

var (
	ErrInvalidArgs  = errors.New("invalid number of arguments")
    ErrInvalidValue = errors.New("invalid argument value")
    ErrStarvation   = errors.New("philosopher starved")
)