package main

import (
	"fmt"
	"time"

	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

type event struct {
	ID      int
	Payload string
}

type ack int

type queue struct {
	id int
	redux.Store
}

func (q *queue) enqueue(prev redux.State, a redux.Action) redux.State {
	prevState := prev.([]event)
	switch t := a.(type) {
	case event:
		newState := make([]event, len(prevState)+1)
		copy(newState, prevState)
		t.ID = q.id
		q.id++
		newState[len(prevState)] = t
		return newState
	case ack:
		var newState []event
		for _, evt := range prevState {
			if evt.ID == int(t) {
				continue
			}
			newState = append(newState, evt)
		}
		return newState
	}
	return prevState
}

func Ack(id int) redux.Thunk {
	return func(dispatch redux.Dispatcher) redux.Action {
		go dispatch(ack(id))
		return id
	}
}

func newQueue() *queue {
	res := new(queue)
	var events []event
	res.Store = redux.New(res.enqueue, events, middleware.Thunk)
	return res
}

func main() {
	q := newQueue()
	q.Subscribe(func() {
		evts := q.GetState().([]event)
		fmt.Println(evts)
		if len(evts) > 0 {
			q.Dispatch(Ack(evts[0].ID))
		}
	})
	for i := 0; i < 100; i++ {
		go q.Dispatch(event{Payload: fmt.Sprint(i)})
	}
	time.Sleep(10 * time.Second)
}
