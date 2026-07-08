package core

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	bus := NewEventBus()

	eventName := "test-event"
	eventData := map[string]string{"foo": "bar"}

	done := make(chan bool)
	bus.On(eventName, func(data interface{}) {
		d := data.(map[string]string)
		if d["foo"] != "bar" {
			t.Errorf("expected bar, got %s", d["foo"])
		}
		done <- true
	})

	bus.Emit(eventName, eventData)

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("callback was not triggered within 1 second")
	}

	// Test PollNativeEvent
	pollResult := bus.PollNativeEvent()
	if pollResult == "" {
		t.Fatal("expected poll result, got empty string")
	}

	var event coreEvent
	if err := json.Unmarshal([]byte(pollResult), &event); err != nil {
		t.Fatalf("failed to unmarshal poll result: %v", err)
	}

	if event.Name != eventName {
		t.Errorf("expected event name %s, got %s", eventName, event.Name)
	}

	// Verify the data in the polled event
	// JSON unmarshaling into interface{} makes it a map[string]interface{}
	polledData := event.Data.(map[string]interface{})
	if polledData["foo"] != "bar" {
		t.Errorf("expected bar in polled data, got %v", polledData["foo"])
	}

	// Test empty poll
	emptyPoll := bus.PollNativeEvent()
	if emptyPoll != "" {
		t.Errorf("expected empty string for empty poll, got %s", emptyPoll)
	}
}
