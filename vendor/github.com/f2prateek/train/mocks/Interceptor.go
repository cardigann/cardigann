package mocks

import "github.com/f2prateek/train"
import "github.com/stretchr/testify/mock"

import "net/http"

func New() *Interceptor {
	return &Interceptor{}
}

type Interceptor struct {
	mock.Mock
}

// Intercept provides a mock function with given fields: _a0
func (_m *Interceptor) Intercept(_a0 train.Chain) (*http.Response, error) {
	ret := _m.Called(_a0)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(train.Chain) *http.Response); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(train.Chain) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
