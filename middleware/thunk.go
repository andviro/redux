package middleware

import "github.com/andviro/redux"

// Thunk allows dispatching action creators
func Thunk(next redux.Dispatcher) redux.Dispatcher {
	return func(a redux.Action) redux.Action {
		switch t := a.(type) {
		case redux.Thunk:
			return t(next)
		}
		return next(a)
	}
}
