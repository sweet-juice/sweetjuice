package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/sweet-juice/sweetjuice/core"
)

type DataStorePlugin struct {
	app *core.Application
}

// NewPlugin creates a new instance of the DataStore plugin.
func NewPlugin() *DataStorePlugin {
	return &DataStorePlugin{}
}

// Name returns the plugin name.
func (p *DataStorePlugin) Name() string {
	return "datastore"
}

// Init initializes the plugin.
func (p *DataStorePlugin) Init(app *core.Application) error {
	p.app = app
	return nil
}

// Set stores a string value for a given key.
func (p *DataStorePlugin) Set(key, value string) error {
	payload, _ := json.Marshal(map[string]string{
		"key":   key,
		"value": value,
	})
	result := core.CallNativePlatform("datastore:set", string(payload))
	return parseResultError(result)
}

// Get retrieves a string value for a given key. Returns defaultVal if key not found.
func (p *DataStorePlugin) Get(key, defaultVal string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"key":     key,
		"default": defaultVal,
	})
	result := core.CallNativePlatform("datastore:get", string(payload))

	if err := parseResultError(result); err != nil {
		return "", err
	}

	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		return "", err
	}
	return resp.Value, nil
}

// Delete removes a key-value pair.
func (p *DataStorePlugin) Delete(key string) error {
	payload, _ := json.Marshal(map[string]string{
		"key": key,
	})
	result := core.CallNativePlatform("datastore:delete", string(payload))
	return parseResultError(result)
}

// Clear removes all data from the store.
func (p *DataStorePlugin) Clear() error {
	result := core.CallNativePlatform("datastore:clear", "{}")
	return parseResultError(result)
}

// GetAll returns all stored key-value pairs.
func (p *DataStorePlugin) GetAll() (map[string]string, error) {
	result := core.CallNativePlatform("datastore:getAll", "{}")
	if err := parseResultError(result); err != nil {
		return nil, err
	}

	var all map[string]string
	if err := json.Unmarshal([]byte(result), &all); err != nil {
		return nil, err
	}
	return all, nil
}

func parseResultError(result string) error {
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(result), &generic); err != nil {
		return nil
	}
	if errMsg, ok := generic["error"]; ok {
		return fmt.Errorf("%v", errMsg)
	}
	return nil
}
