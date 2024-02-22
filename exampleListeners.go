package Nbj

type ExampleListener struct {
}

func (listener *ExampleListener) Handle(event any) bool {
	// Force a specific event by type casting it
	castEvent := event.(*ExampleEvent)

	return !castEvent.StopPropagation
}
