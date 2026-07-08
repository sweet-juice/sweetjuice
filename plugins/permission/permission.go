// Package permission provides a standard plugin for handling Android runtime permissions.
// It bridges the Go application with the native Android permission system.
package permission

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

// PermissionPlugin handles the Go-side logic for the permissions system.
// It allows Go code to check and request permissions via the native bridge.
type PermissionPlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the PermissionPlugin.
func NewPlugin() *PermissionPlugin {
	return &PermissionPlugin{}
}

// Name returns the plugin name "permissions".
func (p *PermissionPlugin) Name() string {
	return "permissions"
}

// Init initializes the plugin with the Sweet Juice application context and registers
// the "permissions:result" native callback handler.
func (p *PermissionPlugin) Init(app *core.Application) error {
	p.app = app

	// Register a handler for the "permissions:result" message coming from Android.
	// This is for Java calling Go.
	app.RegisterNativeMethod("permissions:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var result map[string]interface{}
		if err := json.Unmarshal(args[0], &result); err != nil {
			return nil, err
		}

		// Emit the result as a Sweet Juice event to the frontend.
		app.Events.Emit("permissions:changed", result)
		return map[string]string{"status": "processed"}, nil
	})

	return nil
}

// Check queries the status of a specific permission synchronously.
// Example permission: "android.permission.CAMERA"
func (p *PermissionPlugin) Check(permission string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"permission": permission,
	})

	// Use CallNativePlatform to call into the Android Java bridge.
	// This is Go calling Java.
	result := core.CallNativePlatform("permissions:check", string(payload))
	// convert result into a string status rather than "{status: granted}"
	var resultMap map[string]interface{}
	json.Unmarshal([]byte(result), &resultMap)
	if status, ok := resultMap["status"].(string); ok {
		return status, nil
	}
	return "unknown", nil
}

// Request triggers a native permission request dialog on Android.
// The result will be delivered asynchronously via the "permissions:changed" event.
func (p *PermissionPlugin) Request(permission string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"permission": permission,
	})

	// Use CallNativePlatform to call into the Android Java bridge.
	// This is Go calling Java.
	result := core.CallNativePlatform("permissions:request", string(payload))
	return result, nil
}
