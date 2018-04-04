package main

import (
	"fmt"

	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

type event struct {
	id      int
	payload string
}

type ack int

type sink func(string)

type queue struct {
	id    int
	buf   []event
	sinks []sink
}

type take chan event

func (q queue) Enqueue(payload string) queue {
	buf := make([]event, len(q.buf)+1)
	copy(buf, q.buf)
	q.id++
	e := event{q.id, payload}
	buf[len(q.buf)] = e
	q.buf = buf
	return q
}

func (q queue) Dequeue(id int) queue {
	var buf []event
	for _, evt := range q.buf {
		if evt.id == id {
			continue
		}
		buf = append(buf, evt)
	}
	q.buf = buf
	return q
}

func Emit(event string) redux.Thunk {
	return func(dispatch redux.Dispatcher, getState func() redux.State) redux.Action {
		q := getState().(queue)
		if len(q.sinks) == 0 {
			return dispatch(event)
		}
		for _, sink := range q.sinks {
			sink(event)
		}
		return event
	}
}

func Sink(sink sink) redux.Thunk {
	return func(dispatch redux.Dispatcher, getState func() redux.State) redux.Action {
		dispatch(sink)
		q := getState().(queue)
		for _, ev := range q.buf {
			dispatch(ev.id)
			for _, sink := range q.sinks {
				sink(ev.payload)
			}
		}
		return sink
	}
}

func (q queue) AddSink(t sink) queue {
	sinks := make([]sink, len(q.sinks)+1)
	copy(sinks, q.sinks)
	sinks[len(q.sinks)] = t
	q.sinks = sinks
	return q
}

func newQ() redux.Store {
	res := redux.New(func(prev redux.State, a redux.Action) redux.State {
		q := prev.(queue)
		switch t := a.(type) {
		case sink:
			q = q.AddSink(t)
		case string:
			q = q.Enqueue(t)
		case int:
			q = q.Dequeue(t)
		}
		return q
	}, queue{}, middleware.Thunk)
	return res
}

func main() {
	q := newQ()
	q.Dispatch(Emit("hello"))
	q.Dispatch(Emit("world"))
	q.Dispatch(Emit("i"))
	q.Dispatch(Sink(func(s string) {
		fmt.Println("***", s)
	}))
	q.Dispatch(Emit("am"))
	q.Dispatch(Sink(func(s string) {
		fmt.Println("===", s)
	}))
	q.Dispatch(Emit("buffered"))
	q.Dispatch(Emit("queue"))
	q.Dispatch(Sink(func(s string) {
		fmt.Println("---", s)
	}))
	q.Dispatch(Emit("with"))
	q.Dispatch(Emit("multiple"))
	q.Dispatch(Emit("subscribers"))
}
