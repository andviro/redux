package middleware

import (
	"github.com/andviro/redux"
)

var _ = redux.Middleware(Thunk(redux.Store(nil)))

// Thunk allows dispatching action creators
func Thunk(store redux.GetStateDispatcher) redux.Middleware {
	return func(next redux.Dispatcher) redux.Dispatcher {
		return func(a redux.Action) redux.Action {
			switch t := a.(type) {
			case redux.Thunk:
				return t(store.Dispatch, store.GetState)
			}
			return next(a)
		}
	}
}
