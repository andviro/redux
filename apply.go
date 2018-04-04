package redux

type withDispatch struct {
	Store
	dispatcher Dispatcher
}

func (wd *withDispatch) Dispatch(a Action) Action {
	return wd.dispatcher(a)
}

// ApplyMiddleware returns new store with modified middleware chain
func ApplyMiddleware(store Store, mws ...Middleware) Store {
	if len(mws) == 0 {
		return store
	}
	res := &withDispatch{Store: store, dispatcher: store.Dispatch}
	for i := len(mws) - 1; i >= 0; i-- {
		res.dispatcher = mws[i](res, res.dispatcher)
	}
	return res
}
