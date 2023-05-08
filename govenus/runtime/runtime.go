package runtime

type Task = func(RuntimeContext) bool

type Promise interface {
	IsDone() bool
	OnDone(Task) Promise
	OnDoneWith(Task, RuntimeContextBuilder) Promise
}

type Runtime interface {
	NewContext() RuntimeContextBuilder
	Launch(Task) Promise
	LaunchWith(Task, RuntimeContextBuilder) Promise
	Start()
	Stop()
}
