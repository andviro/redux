package main

import (
	"fmt"
	"time"

	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

type event string

type sink func(event)

type drop struct{}

type queue struct {
	buf   []event
	sinks []sink
}

// Emit sends event to sink or buffer
func Emit(e event) redux.Thunk {
	return func(dispatch redux.Dispatcher, getState func() redux.State) redux.Action {
		dispatch(e)
		q := getState().(queue)
		for _, sink := range q.sinks {
			sink(e)
		}
		return e
	}
}

// Sink subscribes function to event and flushes buffer
func Sink(s sink) redux.Thunk {
	return func(dispatch redux.Dispatcher, getState func() redux.State) redux.Action {
		dispatch(s)
		q := getState().(queue)
		if len(q.sinks) == 1 {
			for _, evt := range q.buf {
				s(evt)
			}
			dispatch(drop{})
		}
		return s
	}
}

// NewQ creates queue based on redux.Store
func NewQ() redux.Store {
	res := redux.New(func(prev redux.State, action redux.Action) redux.State {
		q := prev.(queue)
		switch t := action.(type) {
		case sink:
			q.sinks = append(q.sinks, t)
		case event:
			if len(q.sinks) == 0 {
				q.buf = append(q.buf, t)
			}
		case drop:
			q.buf = nil
		}
		return q
	}, queue{}, middleware.Thunk)
	return res
}

func main() {
	q := NewQ()
	go q.Dispatch(Emit("hello"))
	go q.Dispatch(Emit("world"))
	go q.Dispatch(Emit("i"))
	go q.Dispatch(Sink(func(s event) {
		fmt.Println("***", s)
	}))
	go q.Dispatch(Emit("am"))
	go q.Dispatch(Sink(func(s event) {
		fmt.Println("===", s)
	}))
	go q.Dispatch(Emit("buffered"))
	go q.Dispatch(Emit("event"))
	go q.Dispatch(Emit("queue"))
	go q.Dispatch(Sink(func(s event) {
		fmt.Println("---", s)
	}))
	go q.Dispatch(Emit("with"))
	go q.Dispatch(Emit("multiple"))
	go q.Dispatch(Emit("subscribers"))
	time.Sleep(1 * time.Second)
}
