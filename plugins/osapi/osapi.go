package osapi

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type OsApiPlugin struct {
	app *core.Application
}

type OsInfo struct {
	// Common fields
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	Model         string `json:"model"`

	// Android specific
	SdkInt       int    `json:"sdk_int"`
	Release      string `json:"release"`
	Codename     string `json:"codename"`
	Manufacturer string `json:"manufacturer"`

	// iOS specific
	Name                string `json:"name"`
	LocalizedModel      string `json:"localized_model"`
	IdentifierForVendor string `json:"identifier_for_vendor"`
	IsPhysicalDevice    bool   `json:"is_physical_device"`
}

// NewPlugin creates a new instance of the OsApi plugin.
func NewPlugin() *OsApiPlugin {
	return &OsApiPlugin{}
}

// Name returns the plugin name.
func (p *OsApiPlugin) Name() string {
	return "osapi"
}

// Init initializes the plugin.
func (p *OsApiPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// GetInfo returns information about the Android system.
func (p *OsApiPlugin) GetInfo() (OsInfo, error) {
	var info OsInfo
	result := core.CallNativePlatform("osapi:getInfo", "{}")

	if err := json.Unmarshal([]byte(result), &info); err != nil {
		return info, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return info, nil
}
