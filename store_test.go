package redux_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/andviro/redux"
)

func TestStore(t *testing.T) {
	var stateHistory []int
	reducer := func(s redux.State, a redux.Action) redux.State {
		st := s.(int)
		switch t := a.(type) {
		case int:
			return st + t
		}
		return st
	}
	s := redux.New(reducer, 10)
	cancel := s.Subscribe(func() {
		st := s.GetState().(int)
		t.Logf("%+v", st)
		stateHistory = append(stateHistory, st)
	})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Dispatch(1)
		}()
	}
	wg.Wait()
	cancel()
	cancel = s.Subscribe(func() {
		tst := s.GetState().(int)
		if tst != 120 {
			t.Error("invalid state", tst)
		}
	})
	s.Dispatch(100)
	cancel()
	if !reflect.DeepEqual([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, stateHistory) {
		t.Error("invalid state history", stateHistory)
	}
	s.Dispatch("alal")
}

func TestDispatchInListener(t *testing.T) {
	var store redux.Store
	reducer := func(s redux.State, a redux.Action) (res redux.State) {
		st := s.(int)
		switch t := a.(type) {
		case int:
			return st + t
		}
		return st
	}
	store = redux.New(reducer, 0)
	var cancel func()
	cancel = store.Subscribe(func() {
		go store.Dispatch(1)
		cancel()
	})
	store.Dispatch(1)
}

func TestDispatchInReduce(t *testing.T) {
	var store redux.Store
	reducer := func(s redux.State, a redux.Action) (res redux.State) {
		st := s.(int)
		switch t := a.(type) {
		case int:
			res = st + t // so that result is returned
			go store.Dispatch(2)
			return
		}
		return st
	}
	store = redux.New(reducer, 0)
	finish := make(chan struct{})
	cancel := store.Subscribe(func() {
		if store.GetState().(int) > 1 {
			close(finish)
		}
	})
	defer cancel()
	store.Dispatch(1)
	<-finish
	if store.GetState().(int) != 3 {
		t.Error("invalid state", store.GetState())
	}
}

func TestSubscribeInReduce(t *testing.T) {
	var store redux.Store
	var cancel func()
	var n int
	reducer := func(s redux.State, a redux.Action) (res redux.State) {
		st := s.(int)
		switch t := a.(type) {
		case int:
			cancel = store.Subscribe(func() {
				n++
			})
			return st + t // so that result is returned
		case string:
			cancel()
		}
		return st
	}
	store = redux.New(reducer, 0)
	store.Dispatch(1)
	store.Dispatch("1")
	store.Dispatch("2")
	if n != 1 {
		t.Error("invalid n", n)
	}
}
