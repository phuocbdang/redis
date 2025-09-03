package io_multiplexing

const OpRead Operation = 0
const OpWrite Operation = 1

type Operation uint32

type Event struct {
	Fd int
	Op Operation
}

type IOMultiplexing interface {
	Monitor(event Event) error
	Wait() ([]Event, error)
	Close() error
}
