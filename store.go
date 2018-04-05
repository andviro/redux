package redux

import (
	"sync"
	"sync/atomic"
)

type action struct {
	a    Action
	done chan Action
}

type store struct {
	n         int
	state     atomic.Value
	actions   chan action
	listeners atomic.Value
	reducer   atomic.Value
	lsLock    sync.RWMutex
}

var _ Store = (*store)(nil)

type listeners map[int]Listener

// New creates a Store and initializes it with state and default reducer
func New(reducer Reducer, state State, mws ...MiddlewareFactory) Store {
	res := new(store)
	res.actions = make(chan action)
	res.state.Store(state)
	res.reducer.Store(reducer)
	res.listeners.Store((listeners)(nil))
	go func() {
		defer close(res.actions)
		for action := range res.actions {
			if _, ok := action.a.(stop); ok {
				close(action.done)
				break
			}
			reducer := res.reducer.Load().(Reducer)
			res.state.Store(reducer(res.state.Load(), action.a))
			ls := res.listeners.Load().(listeners)
			for _, l := range ls {
				l()
			}
			action.done <- action.a
		}
	}()
	return ApplyMiddleware(res, mws...)
}

func (s *store) ReplaceReducer(r Reducer) {
	s.reducer.Store(r)
}

func (s *store) GetState() State {
	return s.state.Load()
}

func (s *store) unsub(id int) func() {
	return func() {
		s.lsLock.Lock()
		defer s.lsLock.Unlock()
		ls := s.listeners.Load().(listeners)
		newls := make(listeners)
		for k, v := range ls {
			if k == id {
				continue
			}
			newls[k] = v
		}
		s.listeners.Store(newls)
	}
}

func (s *store) Subscribe(f Listener) UnsubscribeFunc {
	s.lsLock.Lock()
	defer s.lsLock.Unlock()
	ls := s.listeners.Load().(listeners)
	newls := make(listeners)
	for k, v := range ls {
		newls[k] = v
	}
	id := s.n
	s.n++
	newls[id] = f
	s.listeners.Store(newls)
	return s.unsub(id)
}

func (s *store) Dispatch(a Action) Action {
	action := action{a: a, done: make(chan Action)}
	s.actions <- action
	return <-action.done
}
