package Nbj

type ExampleEvent struct {
	StopPropagation bool
}

func (event *ExampleEvent) Accessor() string {
	return "example_event"
}
