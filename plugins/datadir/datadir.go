package datadir

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type DataDirPlugin struct {
	app *core.Application
}

type AppDirs struct {
	Files         string `json:"files"`
	Cache         string `json:"cache"`
	ExternalFiles string `json:"external_files"`
	ExternalCache string `json:"external_cache"`
}

// NewPlugin creates a new instance of the DataDir plugin.
func NewPlugin() *DataDirPlugin {
	return &DataDirPlugin{}
}

// Name returns the plugin name.
func (p *DataDirPlugin) Name() string {
	return "datadir"
}

// Init initializes the plugin.
func (p *DataDirPlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// GetDirs returns the standard application directories.
func (p *DataDirPlugin) GetDirs() (AppDirs, error) {
	var dirs AppDirs
	result := core.CallNativePlatform("datadir:getDirs", "{}")

	if err := json.Unmarshal([]byte(result), &dirs); err != nil {
		return dirs, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return dirs, nil
}
