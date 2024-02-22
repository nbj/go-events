package Nbj

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_event_handler_can_be_created(t *testing.T) {
	// Act
	handler := NewEventHandler()

	// Assert
	assert.Equal(t, "*Nbj.EventHandler", reflect.TypeOf(handler).String())
}

func Test_event_handler_can_have_events_and_listeners_set(t *testing.T) {
	// Arrange
	handler := NewEventHandler()

	events := map[string][]Listener{
		"example-event": {
			&ExampleListener{},
		},
	}

	assert.Equal(t, 0, len(handler.events))

	// Act
	handler.SetEventsAndListeners(events)

	// Assert
	assert.Equal(t, 1, len(handler.events))
}

func Test_event_handler_can_have_events_registered(t *testing.T) {
	// Arrange
	handler := NewEventHandler()
	assert.Equal(t, 0, len(handler.events))

	// Act
	handler.Register("example_event")

	// Assert
	assert.Equal(t, 1, len(handler.events))
}

func Test_event_handler_can_have_listeners_added_to_registered_events(t *testing.T) {
	// Arrange
	handler := NewEventHandler()
	handler.Register("example_event")
	assert.Equal(t, 0, len(handler.events["example_event"]))

	// Act
	handler.AddListener("example_event", &ExampleListener{})

	// Assert
	assert.Equal(t, 1, len(handler.events))
	assert.Equal(t, 1, len(handler.events["example_event"]))
	assert.Equal(t, "*Nbj.ExampleListener", reflect.TypeOf(handler.events["example_event"][0]).String())
}

func Test_event_handler_will_register_event_if_listeners_are_added_to_unregistered_event(t *testing.T) {
	// Arrange
	handler := NewEventHandler()
	assert.Equal(t, 0, len(handler.events))

	// Act
	handler.AddListener("example_event", &ExampleListener{})

	// Assert
	assert.Equal(t, 1, len(handler.events))
	assert.Equal(t, 1, len(handler.events["example_event"]))
	assert.Equal(t, "*Nbj.ExampleListener", reflect.TypeOf(handler.events["example_event"][0]).String())
}

func Test_event_handler_can_dispatch_events(t *testing.T) {
	// Arrange
	handler := NewEventHandler()
	handler.AddListener("example_event", &ExampleListener{})

	// Act
	report := handler.Dispatch(&ExampleEvent{
		StopPropagation: false,
	})

	// Assert
	assert.True(t, report.Registered)
	assert.True(t, report.Completed)
	assert.Equal(t, "example_event", report.Accessor)
}

func Test_event_handler_can_stop_propagation_of_dispatched_events(t *testing.T) {
	// Arrange
	handler := NewEventHandler()
	handler.AddListener("example_event", &ExampleListener{})

	// Act
	report := handler.Dispatch(&ExampleEvent{
		StopPropagation: true,
	})

	// Assert
	assert.True(t, report.Registered)
	assert.False(t, report.Completed)
	assert.True(t, report.ListenerReports[0].StoppedPropagation)
}
