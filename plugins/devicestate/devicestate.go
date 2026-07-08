package devicestate

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type DeviceStatePlugin struct {
	app *core.Application
}

type DeviceState struct {
	BatteryLevel   int              `json:"battery_level"`
	IsCharging     bool             `json:"is_charging"`
	BatteryStatus  string           `json:"battery_status"`
	BatteryHealth  string           `json:"battery_health"`
	Temperature    float64          `json:"temperature"`
	IsLowPowerMode bool             `json:"low_power_mode"`
	Connectivity   ConnectivityInfo `json:"connectivity"`
	Timestamp      int64            `json:"timestamp"`
}

type ConnectivityInfo struct {
	IsConnected bool   `json:"is_connected"`
	NetworkType string `json:"network_type"`
	IsRoaming   bool   `json:"is_roaming"`
	IsUnmetered bool   `json:"is_unmetered"`
}

// NewPlugin creates a new instance of the DeviceState plugin.
func NewPlugin() *DeviceStatePlugin {
	return &DeviceStatePlugin{}
}

// Name returns the plugin name.
func (p *DeviceStatePlugin) Name() string {
	return "devicestate"
}

// Init initializes the plugin and registers the native callback that delivers device state updates.
func (p *DeviceStatePlugin) Init(app *core.Application) error {
	p.app = app

	app.RegisterNativeMethod("devicestate:changed", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided")
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		app.Events.Emit("devicestate:changed", payload)
		return map[string]string{"status": "received"}, nil
	})

	return nil
}

// GetState fetches the current device state from the native Android bridge.
func (p *DeviceStatePlugin) GetState() (DeviceState, error) {
	var state DeviceState
	result := core.CallNativePlatform("devicestate:getState", "{}")

	if err := parsePluginError(result); err != nil {
		return state, err
	}

	if err := json.Unmarshal([]byte(result), &state); err != nil {
		return state, err
	}

	return state, nil
}

// StartMonitoring begins native device state event monitoring.
func (p *DeviceStatePlugin) StartMonitoring() (string, error) {
	return p.callNativePlatform("devicestate:startMonitoring")
}

// StopMonitoring disables native device state event monitoring.
func (p *DeviceStatePlugin) StopMonitoring() (string, error) {
	return p.callNativePlatform("devicestate:stopMonitoring")
}

func (p *DeviceStatePlugin) callNativePlatform(method string) (string, error) {
	result := core.CallNativePlatform(method, "{}")
	if err := parsePluginError(result); err != nil {
		return result, err
	}
	return result, nil
}

func parsePluginError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"]; ok {
		return fmt.Errorf("%v", errMsg)
	}
	return nil
}
