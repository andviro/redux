package main

import (
	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

type event struct {
	id      int
	payload string
}

type ack int

type queue struct {
	id  int
	buf []event
}

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

func newQ() redux.Store {
	return redux.New(func(prev redux.State, a redux.Action) redux.State {
		q := prev.(queue)
		switch t := a.(type) {
		case string:
			return q.Enqueue(t)
		case int:
			return q.Dequeue(t)
		}
		return q
	}, queue{}, middleware.Thunk)
}

func main() {
	q := newQ()
	q.Dispatch("hello")
}
