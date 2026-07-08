package biometric

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type BiometricPlugin struct {
	app *core.Application
}

type AuthOptions struct {
	Title              string `json:"title"`
	Subtitle           string `json:"subtitle"`
	Description        string `json:"description"`
	NegativeButtonText string `json:"negative_button_text"`
}

type AuthResult struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}

type CanAuthResult struct {
	Can    bool   `json:"can_authenticate"`
	Status string `json:"status"`
}

// NewPlugin creates a new instance of the Biometric plugin.
func NewPlugin() *BiometricPlugin {
	return &BiometricPlugin{}
}

// Name returns the plugin name.
func (p *BiometricPlugin) Name() string {
	return "biometric"
}

// Init initializes the plugin and registers the native callback for authentication results.
func (p *BiometricPlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("biometric:result", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var result AuthResult
		if err := json.Unmarshal(args[0], &result); err != nil {
			return nil, err
		}

		app.Events.Emit("biometric:result", result)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// CanAuthenticate checks if the device supports biometric authentication and if it is enrolled.
func (p *BiometricPlugin) CanAuthenticate() (CanAuthResult, error) {
	var res CanAuthResult
	result := core.CallNativePlatform("biometric:canAuthenticate", "{}")

	if err := json.Unmarshal([]byte(result), &res); err != nil {
		return res, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return res, nil
}

// Authenticate triggers the biometric authentication prompt.
// Results are delivered asynchronously via the "biometric:result" event.
func (p *BiometricPlugin) Authenticate(options AuthOptions) (string, error) {
	payload, _ := json.Marshal(options)
	return core.CallNativePlatform("biometric:authenticate", string(payload)), nil
}
