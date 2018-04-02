package main

import (
	"fmt"

	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

type event struct {
	ID      int
	Payload string
	done    chan struct{}
}

func (e *event) Wait() {
	<-e.done
}

type ack int

type queue struct {
	id      int
	current *event
	buf     []event
}

func (q queue) Enqueue(e event) queue {
	buf := make([]event, len(q.buf)+1)
	copy(buf, q.buf)
	e.ID = q.id
	q.id++
	buf[len(q.buf)] = e
	q.buf = buf
	q.current = &e
	return q
}

func (q queue) Dequeue(id int) queue {
	var buf []event
	for _, evt := range q.buf {
		if evt.ID == id {
			close(evt.done)
			continue
		}
		buf = append(buf, evt)
	}
	q.buf = buf
	q.current = nil
	return q
}

func Event(payload string) redux.Thunk {
	return func(dispatch redux.Dispatcher) redux.Action {
		res := event{Payload: payload, done: make(chan struct{})}
		go dispatch(res)
		return res
	}
}

func (e *event) Ack() redux.Thunk {
	return func(dispatch redux.Dispatcher) redux.Action {
		ack := ack(e.ID)
		go dispatch(ack)
		return ack
	}
}

func newQueue() redux.Store {
	return redux.New(func(prev redux.State, a redux.Action) redux.State {
		q := prev.(queue)
		switch t := a.(type) {
		case event:
			return q.Enqueue(t)
		case ack:
			return q.Dequeue(int(t))
		}
		return q
	}, queue{}, middleware.Thunk)
}

func main() {
	q := newQueue()
	n := 0
	q.Subscribe(func() {
		evts := q.GetState().(queue)
		if evts.current != nil {
			fmt.Println(evts.current.ID)
			q.Dispatch(evts.current.Ack())
		}
		//         fmt.Println(evts.buf)
		//         //         x, _ := strconv.Atoi(evts.pending.Payload)
		//         //         n += x
		//         if len(evts.buf) > 0 {
		//             q.Dispatch(evts.buf[0].Ack())
		//         }
	})
	var evts []event
	for i := 0; i < 100; i++ {
		q.Dispatch(Event(fmt.Sprint(i)))
	}
	for _, e := range evts {
		e.Wait()
	}
	fmt.Println(n)
}
