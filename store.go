package redux

import (
	"sync"
	"sync/atomic"
)

const (
	isIdle = iota
	isReducing
	isDispatching
)

type act struct {
	a    Action
	done chan Action
}

type store struct {
	n         int
	state     atomic.Value
	stop      chan struct{}
	events    chan act
	listeners atomic.Value
	reducer   atomic.Value
	lsLock    sync.RWMutex
}

type withDispatch struct {
	Store
	dispatcher Dispatcher
}

func (wd *withDispatch) Dispatch(a Action) Action {
	return wd.dispatcher(a)
}

var _ Store = (*store)(nil)

type listeners map[int]Listener

// ApplyMiddleware returns new store with modified middleware chain
func ApplyMiddleware(store Store, mws ...Middleware) Store {
	if len(mws) == 0 {
		return store
	}
	res := &withDispatch{Store: store, dispatcher: store.Dispatch}
	for i := len(mws) - 1; i >= 0; i-- {
		res.dispatcher = mws[i](res.dispatcher)
	}
	return res
}

// New creates a Store and initializes it with state and default reducer
func New(reducer Reducer, state State, mws ...Middleware) Store {
	res := &store{
		stop:   make(chan struct{}),
		events: make(chan act),
	}
	res.state.Store(state)
	res.reducer.Store(reducer)
	res.listeners.Store((listeners)(nil))
	go func() {
		for {
			select {
			case <-res.stop:
				return
			case action := <-res.events:
				reducer := res.reducer.Load().(Reducer)
				res.state.Store(reducer(res.state.Load(), action.a))
				ls := res.listeners.Load().(listeners)
				for _, l := range ls {
					l()
				}
				action.done <- action.a
			}
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

func (s *store) Dispatch(action Action) Action {
	a := act{a: action, done: make(chan Action)}
	select {
	case <-s.stop:
		return action
	case s.events <- a:
		break
	}
	return <-a.done
}
