package middleware

import (
	"github.com/andviro/redux"
)

var _ = redux.Middleware(Thunk)

// Thunk allows dispatching action creators
func Thunk(store redux.Store, next redux.Dispatcher) redux.Dispatcher {
	return func(a redux.Action) redux.Action {
		switch t := a.(type) {
		case redux.Thunk:
			return t(next, store.GetState)
		}
		return next(a)
	}
}
