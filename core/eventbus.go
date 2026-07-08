package core

import (
	"encoding/json"
	"sync"
)

// coreEvent represents a standardized message format matching core v3 event packets.
type coreEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data,omitempty"`
}

type EventCallback func(data interface{})

type EventBus struct {
	mu          sync.RWMutex
	listeners   map[string][]EventCallback
	nativeQueue chan coreEvent
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners:   make(map[string][]EventCallback),
		nativeQueue: make(chan coreEvent, 100),
	}
}

// On registers a Go routine callback for internal system events.
func (b *EventBus) On(eventName string, callback EventCallback) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[eventName] = append(b.listeners[eventName], callback)
}

// Emit broadcasts an event from Go down to both registered Go routines and the JavaScript frontend.
func (b *EventBus) Emit(eventName string, data interface{}) {
	b.mu.RLock()
	callbacks := b.listeners[eventName]
	b.mu.RUnlock()

	for _, cb := range callbacks {
		go cb(data)
	}

	b.nativeQueue <- coreEvent{
		Name: eventName,
		Data: data,
	}
}

// PollNativeEvent is called by the native Android/iOS wrapper thread to consume events non-blockingly.
func (b *EventBus) PollNativeEvent() string {
	select {
	case event := <-b.nativeQueue:
		bytes, err := json.Marshal(event)
		if err != nil {
			return ""
		}
		return string(bytes)
	default:
		return ""
	}
}
