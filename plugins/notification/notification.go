package notification

import (
	"encoding/json"
	"time"

	"github.com/sweet-juice/sweetjuice/core"
)

// Importance defines the notification priority level.
type Importance string

const (
	ImportanceDefault Importance = "DEFAULT"
	ImportanceHigh    Importance = "HIGH"
	ImportanceLow     Importance = "LOW"
	ImportanceMin     Importance = "MIN"
)

// Notification defines the content and behavior of a system notification.
type Notification struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	ChannelID   string     `json:"channel_id"`
	ChannelName string     `json:"channel_name"`
	Importance  Importance `json:"importance"`
}

// NotificationPlugin handles posting system notifications.
type NotificationPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the NotificationPlugin.
func NewPlugin() *NotificationPlugin {
	return &NotificationPlugin{}
}

// Name returns the plugin name.
func (p *NotificationPlugin) Name() string {
	return "notification"
}

// Init initializes the plugin.
func (p *NotificationPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// Post triggers a system notification.
// If the ID is 0, a unique ID will be generated based on the current timestamp.
func (p *NotificationPlugin) Post(n Notification) (string, error) {
	if n.ID == 0 {
		// Basic way to generate a unique ID if not provided
		n.ID = int(time.Now().Unix() % 100000)
	}
	if n.ChannelID == "" {
		n.ChannelID = "default_channel"
	}
	if n.ChannelName == "" {
		n.ChannelName = "General Notifications"
	}
	if n.Importance == "" {
		n.Importance = ImportanceDefault
	}

	payload, _ := json.Marshal(n)
	return core.CallNativePlatform("notification:post", string(payload)), nil
}

// Cancel removes a specific notification by ID.
func (p *NotificationPlugin) Cancel(id int) string {
	payload, _ := json.Marshal(map[string]int{"id": id})
	return core.CallNativePlatform("notification:cancel", string(payload))
}
