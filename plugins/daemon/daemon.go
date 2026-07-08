package daemon

import (
	"encoding/json"

	"github.com/sweet-juice/sweetjuice/core"
)

// Options defines the foreground service configuration.
type Options struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	ChannelID  string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Importance string `json:"importance"` // HIGH, DEFAULT, LOW
}

// DaemonPlugin manages the application lifecycle as a background/foreground service.
type DaemonPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the Daemon plugin.
func NewPlugin() *DaemonPlugin {
	return &DaemonPlugin{}
}

// Name returns the plugin name.
func (p *DaemonPlugin) Name() string {
	return "daemon"
}

// Init initializes the plugin.
func (p *DaemonPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// Start starts the background/foreground service.
// On Android Oreo+, this starts a Foreground Service with a notification.
func (p *DaemonPlugin) Start(opts Options) (string, error) {
	if opts.ChannelID == "" {
		opts.ChannelID = "daemon_channel"
	}
	if opts.ChannelName == "" {
		opts.ChannelName = "Background Service"
	}
	if opts.Importance == "" {
		opts.Importance = "LOW"
	}

	payload, _ := json.Marshal(opts)
	return core.CallNativePlatform("daemon:start", string(payload)), nil
}

// Stop stops the background/foreground service.
func (p *DaemonPlugin) Stop() (string, error) {
	return core.CallNativePlatform("daemon:stop", "{}"), nil
}
