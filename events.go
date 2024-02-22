package Nbj

import "reflect"

type Event interface {
	Accessor() string
}

type Listener interface {
	Handle(event any) bool
}

type EventDispatchReport struct {
	Registered      bool
	Completed       bool
	Accessor        string
	ListenerReports []ListenerDispatchReport
}

type ListenerDispatchReport struct {
	Instance           string
	Completed          bool
	StoppedPropagation bool
}

type EventHandler struct {
	events map[string][]Listener
}

func NewEventHandler() *EventHandler {
	var instance EventHandler
	instance.events = make(map[string][]Listener)

	return &instance
}

func (handler *EventHandler) SetEventsAndListeners(events map[string][]Listener) *EventHandler {
	handler.events = events

	return handler
}

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
