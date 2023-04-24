package govenus

import (
	"errors"
)

type RuntimeContext interface {
	Runtime() Runtime  // The runtime where the context is being run
	IsAvailable() bool // Determines when the context is available to run
}

type RuntimeContextBuilder interface {
	Build() (RuntimeContext, error)
	SetRuntime(Runtime) RuntimeContextBuilder
	AddAvailabilityCondition(func() bool) RuntimeContextBuilder
}

type simpleContext struct {
	runtime           Runtime
	checkAvailability func() bool
}

func NewContextBuilder() RuntimeContextBuilder {
	return &simpleContext{
		runtime:           nil,
		checkAvailability: func() bool { return true },
	}
}

func (sc *simpleContext) Build() (RuntimeContext, error) {
	if sc.runtime == nil {
		return nil, errors.New("cannot build a context without runtime")
	}
	// clone context to allow reuse the builder
	return &simpleContext{
		runtime:           sc.runtime,
		checkAvailability: sc.checkAvailability,
	}, nil
}

func (sc *simpleContext) Runtime() Runtime {
	return sc.runtime
}

func (sc *simpleContext) SetRuntime(runtime Runtime) RuntimeContextBuilder {
	sc.runtime = runtime
	return sc
}

func (sc *simpleContext) AddAvailabilityCondition(condition func() bool) RuntimeContextBuilder {
	prev := sc.checkAvailability
	sc.checkAvailability = func() bool { return condition() && prev() }
	return sc
}

func (sc *simpleContext) IsAvailable() bool {
	return sc.checkAvailability()
}
