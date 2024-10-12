package safego

import (
	"errors"
	"fmt"
)

var PanicCatchFunc = func(name string, p interface{}) {}

func Recover(name string) {
	if p := recover(); p != nil {
		PanicCatchFunc(name, p)
	}
}

func Call(name string, f func()) (err error) {
	defer func() {
		if p := recover(); p != nil {
			PanicCatchFunc(name, p)
			err = errors.New(fmt.Sprint(p))
		}
	}()
	f()
	return
}

func CallError(name string, f func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			PanicCatchFunc(name, p)
			err = errors.New(fmt.Sprint(p))
		}
	}()
	err = f()
	return
}

func Go(name string, f func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				PanicCatchFunc(name, p)
			}
		}()
		f()
	}()
}
