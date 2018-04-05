// Package redux implements Redux store pattern in a Go way
package redux

// State is an arbitrary Go value
type State interface{}

// Action is matched in reducer using type switch and is arbitrary
type Action interface{}

// Reducer consumes state and action and produces new state
type Reducer func(State, Action) State

// Listener is called on each dispatched action
type Listener func()

// UnsubscribeFunc must be called to remove subscription
type UnsubscribeFunc func()

// Dispatcher receives the action and returns it, possibly modified
type Dispatcher func(Action) Action

// Middleware constructs Dispatcher from another Dispatcher
type Middleware func(Dispatcher) Dispatcher

// GetStateDispatcher represents limited store interface
type GetStateDispatcher interface {
	Dispatch(Action) Action // Send action to modify store
	GetState() State        // Get current state
}

// MiddlewareFactory creates middleware using supplied store
type MiddlewareFactory func(GetStateDispatcher) Middleware

// Thunk conditionally applies dispatcher to action
type Thunk func(Dispatcher, func() State) Action

type stop struct{}

// Stop shuts down store when dispatched
var Stop stop

// Store is a redux store
type Store interface {
	GetStateDispatcher
	Subscribe(Listener) UnsubscribeFunc // Subscribe to store changes
	ReplaceReducer(Reducer)             // Set new reducer
}

//go:generate moq -out mock/store.go -pkg mock . Store
