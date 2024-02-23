package Nbj

import "reflect"

// Event
// Interface that describes an event
type Event interface {
	Accessor() string
}

// Listener
// Interface that describes a listener
type Listener interface {
	Handle(event any) bool
}

// EventDispatchReport
// Structure describing an event dispatch report
type EventDispatchReport struct {
	Registered      bool
	Completed       bool
	Accessor        string
	ListenerReports []ListenerDispatchReport
}

// ListenerDispatchReport
// Structure describing a listener dispatch report
type ListenerDispatchReport struct {
	Instance           string
	Completed          bool
	StoppedPropagation bool
}

// EventHandler
// Structure describing the actual event handler
type EventHandler struct {
	events map[string][]Listener
}

// NewEventHandler
// Named constructor that initializes a new event handler
func NewEventHandler() *EventHandler {
	var instance EventHandler
	instance.events = make(map[string][]Listener)

	return &instance
}

// SetEventsAndListeners
// Overrides all events and associated listener with a complete configuration
func (handler *EventHandler) SetEventsAndListeners(events map[string][]Listener) *EventHandler {
	handler.events = events

	return handler
}

// Register
// Registers an event on the event handler
func (handler *EventHandler) Register(accessor string) *EventHandler {
	// We do not want to override existing events and listener
	// So if the event is already registered we simply return
	if _, ok := handler.events[accessor]; ok {
		return handler
	}

	// Register the event and return
	handler.events[accessor] = []Listener{}

	return handler
}

// AddListener
// Adds a listener to an event. The event is automatically registered if it was not beforehand
func (handler *EventHandler) AddListener(accessor string, listener Listener) *EventHandler {
	// If the event accessor we are adding a listener to does not exist
	// we make sure to register the event accessor
	if _, ok := handler.events[accessor]; !ok {
		handler.Register(accessor)
	}

	// Add the listener to the event accessor
	handler.events[accessor] = append(handler.events[accessor], listener)

	return handler
}

// Dispatch
// Dispatches an event triggering all its listeners. This will generate and
// return an event dispatch report
func (handler *EventHandler) Dispatch(event Event) *EventDispatchReport {
	eventReport := EventDispatchReport{
		Registered: false,
		Completed:  false,
		Accessor:   event.Accessor(),
	}

	if listeners, ok := handler.events[eventReport.Accessor]; ok {
		eventReport.Registered = true
		propagationStopped := false

		for _, listener := range listeners {
			listenerReport := ListenerDispatchReport{
				Instance:           reflect.TypeOf(listener).String(),
				Completed:          false,
				StoppedPropagation: false,
			}

			if propagationStopped {
				eventReport.ListenerReports = append(eventReport.ListenerReports, listenerReport)

				continue
			}

			if !listener.Handle(event) {
				propagationStopped = true
				listenerReport.StoppedPropagation = true
			}

			listenerReport.Completed = true

			eventReport.ListenerReports = append(eventReport.ListenerReports, listenerReport)
		}

		if !propagationStopped {
			eventReport.Completed = true
		}
	}

	return &eventReport
}
