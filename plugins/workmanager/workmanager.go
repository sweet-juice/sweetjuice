package workmanager

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sweet-juice/sweetjuice/core"
)

// NetworkType defines the required network state for a task.
type NetworkType string

const (
	NetworkConnected   NetworkType = "CONNECTED"
	NetworkUnmetered   NetworkType = "UNMETERED"
	NetworkNotRoaming  NetworkType = "NOT_ROAMING"
	NetworkNotRequired NetworkType = "NOT_REQUIRED"
)

// Constraints define the conditions that must be met for a task to run.
type Constraints struct {
	NetworkType           NetworkType `json:"network_type"`
	RequiresCharging      bool        `json:"requires_charging"`
	RequiresDeviceIdle    bool        `json:"requires_device_idle"`
	RequiresBatteryNotLow bool        `json:"requires_battery_not_low"`
	RequiresStorageNotLow bool        `json:"requires_storage_not_low"`
}

// DefaultConstraints returns a standard set of non-restrictive constraints.
func DefaultConstraints() Constraints {
	return Constraints{
		NetworkType: NetworkNotRequired,
	}
}

// TaskFunc is the signature for Go functions that run in the background.
type TaskFunc func() error

// WorkManagerPlugin manages background tasks via the Android WorkManager API.
type WorkManagerPlugin struct {
	app      *core.Application
	registry map[string]TaskFunc
	mu       sync.RWMutex
}

// NewPlugin creates a new instance of the WorkManager plugin.
func NewPlugin() *WorkManagerPlugin {
	return &WorkManagerPlugin{
		registry: make(map[string]TaskFunc),
	}
}

// Name returns the plugin name.
func (p *WorkManagerPlugin) Name() string {
	return "workmanager"
}

// Init initializes the plugin and registers the execution handler for native calls.
func (p *WorkManagerPlugin) Init(app *core.Application) error {
	p.app = app

	// Register the handler that the Android Worker will call
	app.RegisterNativeMethod("workmanager:execute", func(args []json.RawMessage) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("no arguments provided to worker")
		}

		var payload struct {
			TaskKey string `json:"task_key"`
		}
		if err := json.Unmarshal(args[0], &payload); err != nil {
			return nil, err
		}

		p.mu.RLock()
		task, exists := p.registry[payload.TaskKey]
		p.mu.RUnlock()

		if !exists {
			return nil, fmt.Errorf("task '%s' not found in Go registry", payload.TaskKey)
		}

		// Execute the Go task logic
		if err := task(); err != nil {
			return map[string]string{"status": "retry", "error": err.Error()}, nil
		}

		return map[string]string{"status": "success"}, nil
	})

	return nil
}

// RegisterTask binds a task name to a Go function.
func (p *WorkManagerPlugin) RegisterTask(name string, fn TaskFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.registry[name] = fn
}

// EnqueueOneTime schedules a task to run once as soon as possible with optional constraints.
func (p *WorkManagerPlugin) EnqueueOneTime(taskKey string, constraints *Constraints) (string, error) {
	if constraints == nil {
		c := DefaultConstraints()
		constraints = &c
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"task_key":    taskKey,
		"constraints": constraints,
	})
	return core.CallNativePlatform("workmanager:enqueueOneTime", string(payload)), nil
}

// EnqueuePeriodic schedules a task to run every N minutes with optional constraints.
func (p *WorkManagerPlugin) EnqueuePeriodic(taskKey string, intervalMinutes int, constraints *Constraints) (string, error) {
	if constraints == nil {
		c := DefaultConstraints()
		constraints = &c
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"task_key":         taskKey,
		"interval_minutes": intervalMinutes,
		"constraints":      constraints,
	})
	return core.CallNativePlatform("workmanager:enqueuePeriodic", string(payload)), nil
}

// IsEnqueued checks if a task with the given key is currently enqueued or running.
func (p *WorkManagerPlugin) IsEnqueued(taskKey string) (bool, error) {
	payload, _ := json.Marshal(map[string]string{
		"task_key": taskKey,
	})
	result := core.CallNativePlatform("workmanager:isEnqueued", string(payload))

	var resp struct {
		Enqueued bool   `json:"enqueued"`
		Error    string `json:"error,omitempty"`
	}
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		return false, err
	}
	if resp.Error != "" {
		return false, fmt.Errorf("%s", resp.Error)
	}
	return resp.Enqueued, nil
}

// CancelAll stops all scheduled work.
func (p *WorkManagerPlugin) CancelAll() string {
	return core.CallNativePlatform("workmanager:cancelAll", "{}")
}
