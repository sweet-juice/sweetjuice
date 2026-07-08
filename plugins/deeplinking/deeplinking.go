package deeplinking

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sweet-juice/sweetjuice/core"
)

type URLHandler func(url string)

type DeepLinkingPlugin struct {
	app      *core.Application
	handlers []URLHandler
	mu       sync.RWMutex
}

// NewPlugin creates a new instance of the DeepLinking plugin.
func NewPlugin() *DeepLinkingPlugin {
	return &DeepLinkingPlugin{}
}

// Name returns the plugin name.
func (p *DeepLinkingPlugin) Name() string {
	return "deeplinking"
}

// Init initializes the plugin and registers the native callback for deep links.
func (p *DeepLinkingPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("deeplinking:received", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no URL provided")
		}

		var url string
		if err := json.Unmarshal(args[0], &url); err != nil {
			return nil, err
		}

		p.mu.RLock()
		handlers := p.handlers
		p.mu.RUnlock()

		for _, handler := range handlers {
			go handler(url)
		}

		// Also emit a Wails event for the frontend
		app.Events.Emit("deeplinking:received", url)

		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// OnURL registers a handler function that will be called when a deep link is received.
func (p *DeepLinkingPlugin) OnURL(handler URLHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, handler)
}
