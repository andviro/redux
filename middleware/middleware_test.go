package middleware_test

import (
	"testing"

	"github.com/andviro/redux"
	"github.com/andviro/redux/middleware"
)

func TestThunk(t *testing.T) {
	reducer := func(s redux.State, a redux.Action) redux.State {
		st := s.(int)
		switch t := a.(type) {
		case int:
			return st + t
		}
		return st
	}
	increment := func(n int) redux.Thunk {
		return func(dispatch redux.Dispatcher) redux.Action {
			if n > 0 {
				return dispatch(1)
			}
			return dispatch(-1)
		}
	}
	s := redux.New(reducer, 0, middleware.Thunk)
	for i := 0; i < 10; i++ {
		s.Dispatch(increment(100))
	}
	if s.GetState().(int) != 10 {
		t.Error("invalid state", s.GetState())
	}
	for i := 0; i < 10; i++ {
		s.Dispatch(increment(-100))
	}
	if s.GetState().(int) != 0 {
		t.Error("invalid state", s.GetState())
	}
}
