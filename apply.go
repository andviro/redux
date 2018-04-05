package redux

type withDispatch struct {
	Store
	dispatcher Dispatcher
}

func (wd *withDispatch) Dispatch(a Action) Action {
	return wd.dispatcher(a)
}

// ApplyMiddleware returns new store with modified middleware chain
func ApplyMiddleware(store Store, mws ...MiddlewareFactory) Store {
	if len(mws) == 0 {
		return store
	}
	var dispatcher = store.Dispatch
	res := &withDispatch{Store: store, dispatcher: func(Action) Action {
		panic("dispatch while applying middleware")
	}}
	for i := len(mws) - 1; i >= 0; i-- {
		dispatcher = mws[i](res)(dispatcher)
	}
	res.dispatcher = dispatcher
	return res
}
